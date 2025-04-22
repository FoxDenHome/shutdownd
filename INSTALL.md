# Windows

- Copy the following into C:\shutdownd

    - `server.pem` from the server

    - https://github.com/FoxDenHome/shutdownd/releases/download/latest/shutdownd-windows-amd64.exe

- Open a cmd as Administrator and run

```powershell
cd /D C:\shutdownd
.\shutdownd-windows-amd64.exe install
```

- In `services.msc`, make sure the service is running and set to automatic start

# Linux

- Open a shell and run 

```bash
git clone https://github.com/FoxDenHome/shutdownd
```

- Copy `server.pem` from the server to `/etc/shutdownd/server.pem`

- Then to initially install and also to update run

```bash
./directory/of/shutdownd/update.sh
```
