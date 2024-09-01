# TODO
Create an image, so that I do not have to build it all the time on the pi

# How to run on raspberrypi linux/arm64

1. Clone repo

2. Create and setup .env file based on .env.example

3. Build and load docker

```bash
docker buildx build --platform linux/arm64 \
  -t harvest_image:latest \
  --load .
```

4. Setup systemd service

```bash
sudo cp systemd_harvest.service /etc/systemd/system
```

```bash
sudo systemctl daemon-reload
sudo systemctl restart systemd_harvest.service
```

```bash
sudo systemctl status systemd_harvest.service
```
