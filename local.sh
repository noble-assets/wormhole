alias simd=./simapp/build/simd

for arg in "$@"
do
    case $arg in
        -r|--reset)
        rm -rf .duke
        shift
        ;;
    esac
done

if ! [ -f .duke/data/priv_validator_state.json ]; then
  simd init validator --chain-id "duke-1" --home .duke &> /dev/null

  simd keys add validator --home .duke --keyring-backend test &> /dev/null
  simd genesis add-genesis-account validator 1000000ustake --home .duke --keyring-backend test
  simd genesis add-genesis-account noble1cyyzpxplxdzkeea7kwsydadg87357qnah9s9cv 1000000uusdc --home .duke --keyring-backend test

  TEMP=.duke/genesis.json
  touch $TEMP && jq '.app_state.staking.params.bond_denom = "ustake"' .duke/config/genesis.json > $TEMP && mv $TEMP .duke/config/genesis.json
  # Temporarily initialize our Wormhole Chain ID as Sei, to be recognized by the local guardian set.
  touch $TEMP && jq '.app_state.wormhole.config.chain_id = 32' .duke/config/genesis.json > $TEMP && mv $TEMP .duke/config/genesis.json
  touch $TEMP && jq '.app_state.wormhole.config.gov_chain = 1' .duke/config/genesis.json > $TEMP && mv $TEMP .duke/config/genesis.json
  touch $TEMP && jq '.app_state.wormhole.config.gov_address = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQ="' .duke/config/genesis.json > $TEMP && mv $TEMP .duke/config/genesis.json
  touch $TEMP && jq '.app_state.wormhole.guardian_sets = {"0":{"addresses":["vvpCnVfNGLf4pNkaLamrSvBdD74="],"expiration_time":0}}' .duke/config/genesis.json > $TEMP && mv $TEMP .duke/config/genesis.json

  simd genesis gentx validator 1000000ustake --chain-id "duke-1" --home .duke --keyring-backend test &> /dev/null
  simd genesis collect-gentxs --home .duke &> /dev/null

  sed -i '' 's/timeout_commit = "5s"/timeout_commit = "1s"/g' .duke/config/config.toml
fi

simd start --home .duke
