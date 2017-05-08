

<h1>GAPE for Docker containers</h1>  
  
Simple recursive filesystem notifier that writes to local or remote syslog servers.  
Compiled version can be found in the bin directory. It was compiled on Debian Jesse, but will  
probably work on any other distribution. This is a version specific for Docker containers  
and will not work on older versions due to changes in the Docker filesystem layout.  
It will try to detect your currently running containers and figure out the union mountpoints. 
It was tested with the AUFS driver only, when time permits I will try others.  
Tested agains Docker version 17.03.1-ce, build c6d412e   
It currently logs any notifications that are known to work  
on any operating system like:  

 - Remove   
 - Create   
 - Write   
 - Rename  


Configuration:  

 - syslogproto: => udp or tcp  
 - sysloghost: Syslog server to send log messages to  
 - syslogport: Syslog server port  
 - localonly: Log only to local syslog (no network) => true or false  
 - stdout: Also output to stout => true or false  
 - paths: An array of paths to watch, these should be directories within the containers 

Usage:

 Download or clone this repo and in the bin directory you can find an example  
 config file, adjust this and copy it to /etc or somewhere else in which case  
 you start it with <pre><code>./gape-docker -config \<path to your config\></code></pre>  
 If you take the default /etc directory simply do <pre><code>./gape-docker</code></pre> or with whatever init or  
 systemd script you want to use.  There are several tools for daemonizing, at a later  
 stage I might build this in (if time permits :)  
 
 Cli output example:   
<pre>
root@debian1:~/goprojects/src/grid6/gape-docker# ./gape-docker
/root/dockerdata/aufs/mnt/65d283960f27d59243fc6f141e83b977e77219b4afb3dc1a28893fd04d173b38/etc
/root/dockerdata/aufs/mnt/65d283960f27d59243fc6f141e83b977e77219b4afb3dc1a28893fd04d173b38/tmp
/root/dockerdata/aufs/mnt/c55d2d99e6fe6f72902b0a1dc157b94f4fad06425464e8122e279e96de6be820/etc
/root/dockerdata/aufs/mnt/c55d2d99e6fe6f72902b0a1dc157b94f4fad06425464e8122e279e96de6be820/tmp
New container detected /root/dockerdata/aufs/mnt/58ed8175b58ecc7577744293788c7bf6a5813b6ce486920e4357405959925e73-init
notify.Write: "/root/dockerdata/aufs/mnt/c55d2d99e6fe6f72902b0a1dc157b94f4fad06425464e8122e279e96de6be820/tmp/kafka-logs/replication-offset-checkpoint.tmp"
New container detected /root/dockerdata/aufs/mnt/58ed8175b58ecc7577744293788c7bf6a5813b6ce486920e4357405959925e73
</pre>   

 If you found any bugs or are using it to full hapiness drop me an <a href="mailto:floris.meester@gmail.com?Subject=GAPE" >email.</a>  

 Based on the excellent library from https://github.com/rjeczalik/notify  
  
  
  
