
ssh pi "sudo systemctl stop jellyfin_uploader.service"
scp bin/jellyfin_uploader jellyfin@pi:/jellyfin/jellyfin_uploader/jellyfin_uploader
scp -r webapp/dist jellyfin@pi:/jellyfin/jellyfin_uploader/dist
