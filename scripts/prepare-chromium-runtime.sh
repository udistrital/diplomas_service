#!/bin/sh
# Ejecutar en el mismo Debian 12 usado por el paso build de Drone.
set -eu

runtime_dir=${1:-runtime}
if [ "$(dpkg --print-architecture)" != amd64 ]; then
    echo "El runtime requiere Debian 12 amd64, igual que la imagen final." >&2
    exit 1
fi
test -x /usr/lib/chromium/chromium
mkdir -p "$runtime_dir/usr/lib" "$runtime_dir/usr/bin" "$runtime_dir/usr/share" "$runtime_dir/etc"
cp -a /usr/lib/chromium "$runtime_dir/usr/lib/"

# Chromium carga NSS dinámicamente; incluir módulos aunque ldd no los liste.
mkdir -p "$runtime_dir/usr/lib/x86_64-linux-gnu"
for library in /usr/lib/x86_64-linux-gnu/libnss*.so /usr/lib/x86_64-linux-gnu/libsoftokn3.so /usr/lib/x86_64-linux-gnu/libfreebl*.so; do
    test -f "$library" || continue
    cp -L "$library" "$runtime_dir/usr/lib/x86_64-linux-gnu/"
done

# Resolver bibliotecas de Chromium y de sus módulos gráficos/NSS.
find "$runtime_dir/usr/lib" -type f \( -name chromium -o -name '*.so' -o -name '*.so.*' \) \
    -exec ldd {} \; 2>/dev/null > "$runtime_dir/ldd.list"
if awk '/not found/ {missing=1; print > "/dev/stderr"} END {exit !missing}' "$runtime_dir/ldd.list"; then
    exit 1
fi
awk '/=> \// {print $3} /^\// {print $1} /^[[:space:]]*\/.*ld-linux/ {print $1}' "$runtime_dir/ldd.list" \
    | sort -u > "$runtime_dir/libraries.list"
while IFS= read -r library; do
    cp -L --parents "$library" "$runtime_dir/"
done < "$runtime_dir/libraries.list"
rm "$runtime_dir/libraries.list" "$runtime_dir/ldd.list"

cp -a /etc/fonts "$runtime_dir/etc/"
cp -a /usr/share/fonts /usr/share/fontconfig "$runtime_dir/usr/share/"
# Evitar wrapper /usr/bin/chromium, que requiere shell en imágenes convencionales.
ln -sf /usr/lib/chromium/chromium "$runtime_dir/usr/bin/chromium"
