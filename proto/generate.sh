cd proto
buf generate --template buf.gen.gogo.yaml
buf generate --template buf.gen.pulsar.yaml
cd ..

cp -r github.com/noble-assets/wormhole/* ./
cp -r api/wormhole/* api/
find api/ -type f -name "*.go" -exec sed -i 's|github.com/noble-assets/wormhole/api/wormhole|github.com/noble-assets/wormhole/api|g' {} +

rm -rf github.com
rm -rf api/wormhole
rm -rf wormhole
