English | [简体中文](translations/README-cn.md)

## 🚀HttpCat Overview
HttpCat is an HTTP file transfer service, designed to provide a simple, efficient, and stable solution for file uploading and downloading.

Project goals: To create a reliable, efficient, and user-friendly HTTP file transfer Swiss Army Knife that greatly enhances your control and experience with file transfers.

Whether it's for temporary sharing or bulk file transfers, HttpCat will be your excellent assistant.

Please note that this translation is a direct translation and may require further refinement by a professional translator for the best results.

## 💥Key Features
* Simple and easy to use
* No external dependencies, easy to port

## 🎉Installation
### Quick Installation
You can directly download the latest httpcat installation package.

After extracting it, simply run the `install.sh` script to install.

```bash
httpcat_version="v0.1.5"
mkdir target_directory
tar -zxvf httpcat_$httpcat_version.tar.gz -C target_directory
```

```bash
cd target_directory/release
./install.sh
```

```bash
systemctl status httpcat
systemctl stop httpcat
systemctl start httpcat

tail -f /root/log/httpcat.log
```


#### For versions prior to v0.1.2
1. Download the latest httpcat installation package.
   `https://github.com/shepf/httpcat-release/tags`

2. Planning and creating directories for httpcat usage
   Assuming we plan to start the project as follows:
   ```bash
   /usr/local/bin/httpcat  --port=80 --static=/home/web/website/httpcat_web/  --upload=/home/web/website/upload/ --download=/home/web/website/upload/  -C /etc/httpdcat/svr.yml
   ```
   * --port Specify the listening port for httpcat.
   * --upload Specify the directory for uploading files.
   * --download Specify the directory for downloading files.
   * -C Specify the configuration file to use. (Note: Modify the location of the SQLite file storage as needed: sqlite_db_path: "./data/sqlite.db")


Prepare the directory for file uploads (we will use the same directory for uploads and downloads):
   ```bash
   mkdir -p /home/web/website/upload/
   ```

Prepare the web static resource directory.
   ```bash
   mkdir -p /home/web/website/httpcat_web/  
   ```

Prepare the directory for storing configuration files.
   ```bash
   mkdir -p /etc/httpdcat/
   ```

3. Installation
```bash
   mkdir httpcat
   cd httpcat
```
Upload the installation package: httpcat_v0.1.1.tar.gz、httpcat_web_v0.1.1.zip


install httpcat
```bash
tar -zxvf httpcat_v0.1.1.tar.gz
cp httpcat /usr/local/bin/
cp conf/svr.yml /etc/httpdcat/
```

install httpcat_web
```bash
cp httpcat_web_v0.1.1.zip /home/web/website/
cd /home/web/website/
unzip httpcat_v0.1.1.tar.gz
mv dist httpcat_web
```

check
```bash
httpcat -v
httpcat -h
```

The command-line parameters for running on Windows are the same as on Linux, except that you use httpcat.exe instead of httpcat.
```bash
httpcat.exe --upload /home/web/website/download/ --download /home/web/website/download/ -C F:\open_code\httpcat\server\conf\svr.yml
```

### Run in the background using tmux
You can use tmux to run in the background:
```bash
Create a new tmux session using a socket file named tmux_httpcat
$ tmux -S tmux_httpcat

#  Once inside tmux, you can execute running commands, such as:
httpcat --static=/home/web/website/upload/  -C server/conf/svr.yml

Move process to background by detaching
Ctrl+b d OR ⌘+b d (Mac)

To re-attach
$ tmux -S tmux_httpcat attach

Alternatively, you can use the following single command to both create (if not exists already) and attach to a session:
$ tmux new-session -A -D -s tmux_httpcat

To delete farming session
$ tmux kill-session -t tmux_httpcat
```


### Linux can use systemd to run in the background
The installation package comes with an httpcat.service file that allows you to run httpcat in the background using systemd. 

You can modify the httpcat.service file according to your needs.

For example, you can modify the ExecStart parameter in the httpcat.service file to specify your own startup parameters.

To add a listening port parameter, you can add --port=80:
```bash
ExecStart=/usr/local/bin/httpcat --port=80  --static=/home/web/website/httpcat_web/  --upload=/home/web/website/upload/ --download=/home/web/website/upload/  -C /etc/httpdcat/svr.yml
```

```bash
cp  httpcat.service /usr/lib/systemd/system/httpcat.service
sudo systemctl daemon-reload
sudo systemctl start httpcat
```

> Note: You may need to modify the startup parameters according to your needs.
> Ensure that the following three directories are consistent (so that the upload directory is also the download directory,
> and it is also the web frontend directory where files can be downloaded without authentication).

```bash
vi httpcat.service
```
```
ExecStart=/usr/local/bin/httpcat  --static=/home/web/website/upload/  --upload=/home/web/website/upload/ --download=/home/web/website/upload/  -C /etc/httpdcat/svr.yml
```

## httpcat web frontend
The new version includes a frontend page. Prior to v0.1.1, the frontend was released separately, and users could choose to download it according to their needs. 

Starting from v0.1.2, the frontend is directly integrated into the installation package, eliminating the need to separately download frontend files.

Since httpcat comes with built-in static resource file handling, users have the freedom to decide whether to use the frontend page.

This frontend is a single-page application. In the production environment, static resources are accessed through the /static route, while API endpoints are accessed through the /api route. 

If users set up their own Nginx server, they should configure the /static route to point to the static resource directory and the /api route to the httpcat service.

To use the frontend, download the release package and extract it to the web directory. httpcat will automatically load the static resource files from the web directory.

The web directory is specified in the configuration file using the static parameter. If not specified, the default location is the website/static directory under the current directory.

Alternatively, you can specify the directory using command-line parameters, such as:
```bash
--static=/home/web/website/httpcat_web/
```

### Frontend Deployment
1. Download the standalone frontend release file, such as httpcat_web_xxx.zip.
2. Extract it to the web directory
    ```bash
       cd /home/web/website/
       unzip httpcat_web_v0.1.1.zip
       mv  dist httpcat_web
    ```
3. Starting the httpcat Service
   To start the service, you need to specify the web interface directory using the --static parameter. For example:
    ```bash
    ./httpcat --static=/home/web/website/httpcat_web/  -C conf/svr.yml
    ```
4. Accessing the httpcat Frontend Service
    ```bash
    http://127.0.0.1:8888
    ```


## ❤ Tips and Tricks
### File Operation Related APIs
#### Uploading Files Using Curl Tool
```bash
curl -v -F "f1=@/root/hello.mojo" -H "UploadToken: httpcat:dZE8NVvimYNbV-YpJ9EFMKg3YaM=:eyJkZWFkbGluZSI6MH0=" http://localhost:8888/api/v1/file/upload
```
The curl command is used to send a multipart/form-data format POST request to the specified URL. Here is an explanation of each part:
- `curl`:  curl is a tool used for transferring data to/from a server, supporting multiple protocols.
- `-v`: Detailed operational information is displayed during command execution, which is known as verbose mode.
- `-F "f1=@/root/hello.mojo"`: Specifies the form data to be sent. The -F option indicates that a form is being sent, and f1=@/root/hello.mojo indicates that the file field to be uploaded is named f1, with the file path /root/hello.mojo. The value of this field is the relative or absolute path to the local file.
- `http://localhost:8888/api/v1/file/upload`: 要发送请求到的 URL，这条命令会将文件上传到这个 URL。
- `-H "UploadToken: httpcat:dZE8NVvimYNbV-YpJ9EFMKg3YaM=:eyJkZWFkbGluZSI6MH0="`: "Upload Token" is a unique authentication token generated based on the "app_key" and "app_secret". When uploading a file, the token is attached and the server verifies its validity.

> Note: f1 is defined in the server-side code. Modifying it to something else, such as file, will result in an error and the upload will fail.

In the curl command, you can use the --retry parameter to specify the number of retry attempts after a failure.

By setting the --retry parameter to a value greater than 0, you can instruct curl to retry uploading the file if it fails.
```
curl --retry 3 xxx
```
You can adjust the retry count based on the actual situation to ensure the reliability and stability of file uploads.


#### File Upload Authentication: UploadToken
If the configuration file has enable_upload_token enabled, file uploads require authentication.

You need to add the upload token to the request header, with the token value being the same as the enable_upload_token value in the configuration file.

An independent upload token credential is generated based on the app_key and app_secret. 

When uploading a file, the upload token is included, and the server will verify the token's validity.

Upload token is generated based on app_key and app_secret. The system will have a built-in app_key and app_secret according to the configuration file.

> Note: The built-in app_key and app_secret in the system can only be modified through the svr.yml file and cannot be modified through the interface.
> Restarting httpcat will load the system's built-in app_key and app_secret.

svr.yml：
```bash
app_key: "httpcat"
app_secret: "httpcat_app_secret"
```

In addition to the built-in app_key and app_secret, you can also add custom app_key and app_secret through the interface. 
You can generate upload tokens based on the app_key and app_secret through the interface.
As shown in the figure below, you can click the "Generate Upload Token" button to obtain the upload token.
![img.png](translations/img.png)


####  Upload file Enterprise WeChat webhook notification.
Configure the persistent_notify_url in the svr.yml file. After a successful upload, an Enterprise WeChat notification will be sent.

The notification message is as follows:

File upload archived, upload information:
- IP Address: 192.168.31.3
- Upload Time: 2023-11-29 23:07:04
- File Name: syslog.md
- File Size: 4.88 KB
- File MD5: 8346ecb8e6342d98a9738c5409xxx

#### Support SQLite to retain upload history.
If the enable_sqlite option is enabled in the configuration, uploaded files will be recorded in an SQLite database. 

You can use the `sqlite` command-line tool to query the upload history records.

Use the `sqlite` command-line tool to create a database and query data.

```bash
sudo apt install sqlite3
sqlite3 --version
```

Run the following command to connect to an SQLite database and specify the filename for the database to be created (e.g., sqlite.db):
```bash
sqlite3 sqlite.db
```

At the sqlite3 prompt, enter the .tables command to list all the tables in the database:
```bash
.tables
```

```bash
SELECT * FROM notifications;
```


#### download file
Download a specific file.
```bash
wget -O xxx.jpg  http://127.0.0.1:8888/api/v1/file/download?filename=xxx.jpg
```
When using the wget command to download a file, the name of the file is determined by the filename portion of the request URL.

Due to the presence of URL parameters, the wget command may treat the entire URL as the filename.

To ensure the correct filename for the downloaded file, you can use the -O parameter to specify the filename.

### P2P Related APIs
P2P functionality needs to be enabled in the configuration file, which is disabled by default.

#### Sending messages to the P2P network via HTTP API
```bash
http://{{ip}}:{{port}}/api/v1/p2p/send_message
POST
{
"topic": "httpcat",
"message": "ceshi cccccccccccc"
}
```

## 💪TODO
1. HTTPS support

Feel free to raise an issue. Good luck! 🍀

## 🍀 FAQ
### What to Do If You Forget the Password?
If you forget the password, you can modify the SQLite database by deleting the admin user. After restarting the httpcat service, a new admin user will be created.

Alternatively, you can directly delete the SQLite database and restart the httpcat service, which will create a new SQLite database.

The default path for the SQLite database is specified by the sqlite_db_path parameter in the configuration file, which is set to ./data/sqlite.db by default. You can modify the SQLite database path by changing this configuration.

To reset the password, locate the SQLite database file and delete it. Then, restart the httpcat service, and a new SQLite database will be created.
```bash
find / -name httpcat_sqlite.db
rm /root/data/httpcat_sqlite.db

systemctl status httpcat
systemctl stop httpcat
systemctl start httpcat
```

## 📝License
1. This software is for personal use only and is strictly prohibited for commercial purposes.
2. The copying, distribution, modification, and use of this software is subject to the following conditions:
   - Prohibited for commercial purposes.
   - Prohibited for use in any commercial product or service.
   - Copyright and license statements must be preserved within the software.
   - Modifying or removing copyright and license statements within the software is prohibited unless explicitly permitted.
3. This software is provided "as is" without any warranties or liabilities.
4. By using this software, you indicate that you have accepted this license agreement.

Welcome to follow our GitHub project, ✨star it to stay updated with our real-time developments. Good luck! 🍀

