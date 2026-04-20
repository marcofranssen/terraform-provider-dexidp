#!/bin/bash
VERSION="1.0.0"
PROVIDER_NAME="dexidp"
PLATFORMS=("linux_amd64" "windows_amd64")

# Bepaal de locatie van de bestanden
if [ -d "dist" ]; then
    PREFIX="dist/"
else
    PREFIX=""
fi

# Start de opbouw van de JSON
JSON_CONTENT="{\n  \"archives\": {"

for i in "${!PLATFORMS[@]}"; do
    OS_ARCH="${PLATFORMS[$i]}"
    ZIP_FILE=$(ls ${PREFIX}terraform-provider-${PROVIDER_NAME}_${VERSION}_${OS_ARCH}.zip 2>/dev/null)

    if [ -f "$ZIP_FILE" ]; then
        # Bereken de h1 hash
        HASH=$(openssl dgst -sha256 -binary "$ZIP_FILE" | openssl base64 | sed 's/^/h1:/')
        FILENAME=$(basename "$ZIP_FILE")

        # Voeg platform data toe aan de string
        JSON_CONTENT+="\n    \"$OS_ARCH\": {\n      \"hashes\": [\"$HASH\"],\n      \"url\": \"$FILENAME\"\n    }"
        
        # Voeg een komma toe als het niet de laatste is
        if [ $i -lt $((${#PLATFORMS[@]} - 1)) ]; then
            JSON_CONTENT+=","
        fi
        echo "✅ $OS_ARCH toegevoegd."
    else
        echo "⚠️  Waarschuwing: Bestand voor $OS_ARCH niet gevonden ($ZIP_FILE)."
    fi
done

# Sluit de JSON af
JSON_CONTENT+="\n  }\n}"

# Schrijf naar bestand
echo -e "$JSON_CONTENT" > "${VERSION}.json"
echo "🚀 ${VERSION}.json is nu compleet."

# Update index.json
echo "{\"versions\": {\"${VERSION}\": {}}}" > index.json