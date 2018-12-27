procDir="./proc/"


uuid=$1
inDir="$procDir$uuid"
outDir="$
echo $inDir
find $inDir -maaxdepth 1 -iname "*.jpg" | xargs -L1 -I{} convert -resize 30% "{}" _resized/"{}"
