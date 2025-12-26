LATEST_TAG=$(git describe --tags --abbrev=0)

IFS='.' read -r MAJOR MINOR PATCH <<< "${LATEST_TAG#v}"

if (( PATCH < 9 )); then
  ((PATCH++))
else
  ((MAJOR++))
  PATCH=0
fi

NEW_TAG="v${MAJOR}.${MINOR}.${PATCH}"

echo "Current tag: $LATEST_TAG"
echo "New tag:     $NEW_TAG"

git add .
git commit -m "building release for tag $NEW_TAG"
git push origin $NEW_TAG
