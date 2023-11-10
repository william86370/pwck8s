sudo docker run -d --restart=unless-stopped -p 443:443 -p 6445:6445 --privileged rancher/rancher:v2.7.1

docker logs  e647f168b62812cd01e33ddf2c8aeb14026661b2f04e2166a7e7525102ef6f66  2>&1 | grep "Bootstrap Password:"