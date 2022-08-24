#!/bin/bash

# Run HID-Node Chain
echo "Setting up hid-node chain"
echo ""
../../../scripts/localnet-single-node/setup.sh
echo ""
echo "Setup Done"
echo ""

echo "Running hid-node"
echo ""
tmux new -s hidnode -d hid-noded start
sleep 5
if [[ -n $(hid-noded status) ]]; then
  echo "hid-noded daemon is now running"
  echo ""
else
  echo "hid-noded daemon failed to start, exiting...."
  exit 1
fi

# Run Osmosis Chain
echo "Setting up Osmosis Chain"
echo ""
./osmosis/osmosis_setup.sh
echo ""
echo "Setup Done"
echo ""

echo "Running osmosis"
echo ""
tmux new -s osmosisnode -d osmosisd start
sleep 5
osmosisd status &> /dev/null
RET_VAL=$?
if [ ${RET_VAL} -eq 0 ]; then
  echo "osmosisd daemon is now running"
  echo ""
else
  echo "osmosisd daemon failed to start, exiting...."
  exit 1
fi

# Run Hermes Relayer
echo "Setting up hermes relayer"
HID_NODE_VALIDATOR_WALLET=$(hid-noded keys show node1 -a)
OSMOSIS_VALIDATOR_WALLET=$(osmosisd keys show osmonode1 -a)
./hermes/setup.sh ${HID_NODE_VALIDATOR_WALLET} ${OSMOSIS_VALIDATOR_WALLET}
echo ""
sleep 3
echo "Starting hermes relayer"
tmux new -s hermesrelayer -d hermes start
sleep 2
echo "Hermes has been started"

echo "Transferring tokens from HID Node to Osmosis"
IBC_TRANSFER_RESULT=$(hid-noded tx ibc-transfer transfer transfer channel-0 ${OSMOSIS_VALIDATOR_WALLET} 1234uhid --broadcast-mode block --from ${HID_NODE_VALIDATOR_WALLET} --output json --yes)

CODE=$(echo ${IBC_TRANSFER_RESULT} | jq '.code')
if [ ${CODE} -eq 0 ]; then
  echo "Token is transferred successfully through IBC"
else
  echo "HID Token did not went through IBC"
  exit 1
fi