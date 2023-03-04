# devip - spawn IPs locally for testing 
This is an experiment with Chat GPT, most of the code has been written by it following my instructions and RFCs. The goal was to make a tool that I can use to quickly "fake" IPs locally for test scenarios, i.e. have my local machine respond to requests to that IP. After some convoluted and not quite working pure-go implementations I allowed it to use `ip` and `ping` and got something useful out of it. 

# How To Use
First of all, build and install it:
```bash
sudo CGO_ENABLED=0 go build -o  /usr/local/bin/devip  .
```

## Add IPs
```
devip add [IP 1] <IP 2> .. <IP n>
```
For example:
```bash
devip add 10.10.10.1 20.20.20.1 30.30.30.1
```
This will make `localhost` serve the IPs `10.10.10.1`, `20.20.20.1` and `30.30.30.1` by creating loopback aliases for them, so you can use them for development. Be aware that this means you won't be able to connect to the real IPs while the aliases are up.

## Remove IPs
```
devip remove [IP 1] <IP 2> .. <IP n>
```
For example:
```bash
devip remove 10.10.10.1 20.20.20.1 30.30.30.1
```
This will stop `localhost` serving the IPs `10.10.10.1`, `20.20.20.1` and `30.30.30.1` by removing the loopback aliases. 

## List IPs
```bash
devip 
```
This will print all IPs currently served by `localhost`. 

## Nginx Example
Let's say I want to start a bunch of services with HTTP frontends without changing ports for each one. Instead, I can just give each its own IP:
```bash
# add the IPs to use
devip add 192.168.0.1 
devip add 10.0.0.5 
devip add 56.65.56.65 

# let's create the web roots
sudo mkdir -p /var/www/html/192.168.0.1/
sudo mkdir -p /var/www/html/10.0.0.5/
sudo mkdir -p /var/www/html/56.65.56.65/

# set permissions so we can write
sudo chmod 0777 /var/www/html/192.168.0.1/
sudo chmod 0777 /var/www/html/10.0.0.5/
sudo chmod 0777 /var/www/html/56.65.56.65/

# create some test files
echo "192.168.0.1" > /var/www/html/192.168.0.1/index.html
echo "10.0.0.5" > /var/www/html/10.0.0.5/index.html
echo "56.65.56.65" > /var/www/html/56.65.56.65/index.html

# set permissions to sensible defaults
sudo chmod 0755 /var/www/html/192.168.0.1/
sudo chmod 0755 /var/www/html/10.0.0.5/
sudo chmod 0755 /var/www/html/56.65.56.65/

# own the stuff
sudo chown -R www-data:www-data /var/www/html/192.168.0.1
sudo chown -R www-data:www-data /var/www/html/10.0.0.5
sudo chown -R www-data:www-data /var/www/html/56.65.56.65

# let's edit the nginx config
sudo nano /etc/nginx/sites-enabled/default
```
```nginx
# for this example we'll use this config:

server { 
	listen 192.168.0.1:80; 
	root /var/www/html/192.168.0.1/; 
	index index.html; 
	server_name machine-a; 
	location / { 
		try_files $uri $uri/ =404;
    }
}

server { 
	listen 10.0.0.5:80; 
	root /var/www/html/10.0.0.5/; 
	index index.html; 
	server_name machine-b; 
	location / { 
		try_files $uri $uri/ =404;
    }
}

server { 
	listen 56.65.56.65:80; 
	root /var/www/html/56.65.56.65/; 
	index index.html; 
	server_name machine-c; 
	location / { 
		try_files $uri $uri/ =404;
    }
}
```

```bash
# restart nginx
sudo service nginx restart

# curl the IPs to confirm it works
MACHINE_A=$(curl -s http://192.168.0.1) 
MACHINE_B=$(curl -s http://10.0.0.5) 
MACHINE_C=$(curl -s http://56.65.56.65)

# print the results
echo "Machine A (192.168.0.1) identifies as $MACHINE_A"
echo "Machine B (10.0.0.5)    identifies as $MACHINE_B"
echo "Machine C (56.65.56.65) identifies as $MACHINE_C"
```

And when done using the IPs you can remove them:
```bash
# remove the IPs we used
devip remove 192.168.0.1 
devip remove 10.0.0.5 
devip remove 56.65.56.65 
```
