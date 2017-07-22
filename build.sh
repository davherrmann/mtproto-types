echo "Downloading latest MTProto scheme..."
wget -O scheme.tl -q https://raw.githubusercontent.com/telegramdesktop/tdesktop/master/Telegram/Resources/scheme.tl
echo "Generating Go file..."
./generate.sh > types.go
