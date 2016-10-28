sunc
====

sunc - Simple Unprivileged Container.


Usage
-----

Start minimal test container.

    $ cd $(mktemp -d)
    $ cp -a /bin /lib /lib64 .
    $ sunc
    # 


### Using the docker container's filesystem

First, prepare the tar archive.

    $ docker export $(docker create nginx:alpine) > nginx-alpine.tar

And start container.
    
    $ cd $(mktemp -d) && chmod +x .
    $ tar xf nginx-alpine.tar
    $ sed -i 's/\(^\s*listen.\+\)80;$/\18080;/' etc/nginx/conf.d/default.conf
    $ sunc usr/sbin/nginx -g "daemon off;"


Install
-------

    make && sudo make install

