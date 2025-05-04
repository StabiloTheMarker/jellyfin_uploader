
Set-Location webapp
npm run build
Set-Location ..
ssh pi "sudo systemctl stop jellyfin_uploader.service"
ssh jellyfin@pi "rm -r /jellyfin/src"
ssh jellyfin@pi "mkdir /jellyfin/src"
scp -r handlers jellyfin@pi:/jellyfin/src
scp -r models jellyfin@pi:/jellyfin/src
scp -r repositories jellyfin@pi:/jellyfin/src
scp -r util jellyfin@pi:/jellyfin/src
scp main.go jellyfin@pi:/jellyfin/src
scp -r webapp/dist jellyfin@pi/jellyfin
ssh jellyfin@pi "cd /jellyfin/src && go mod init jellyfin_uploader"
ssh jellyfin@pi "cd /jellyfin/src && go get -u gorm.io/gorm"
ssh jellyfin@pi "cd /jellyfin/src && go get -u gorm.io/driver/sqlite"
ssh jellyfin@pi "cd /jellyfin/src && go mod tidy"
ssh jellyfin@pi "cd /jellyfin/src  && go build -o '../jellyfin_uploader' ."
ssh pi "sudo systemctl restart jellyfin_uploader.service"
