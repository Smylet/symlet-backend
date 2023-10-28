#!/bin/bash

# Encode app_test.env.sample to base64 and update GitHub secret


# Navigate to .git/hooks and create a file called pre-push
# Add the following code to the file
# Run chmod +x .git/hooks/pre-push to make the file executable
# run gh auth login and add your Github Token(classic)
# Run ./.git/hooks/pre-push to execute the file 


#!/bin/bash

base64 -i $(pwd)/resources/env/app_test.env -o gh.txt

# Update the GitHub secret using the GitHub CLI

gh secret set APP_ENV_TEST -R https://github.com/Smylet/symlet-backend  -b "$(cat gh.txt)"

exit 0





