# -*- mode: ruby -*-
# vi: set ft=ruby :

GO_VERSION="1.7.1"

Vagrant.configure("2") do |config|
  config.vm.box = "boxcutter/ubuntu1604"

  config.vm.provision "shell", inline: <<-SHELL
    apt update
    apt install -y git

    if [ ! -e /usr/local/go ]; then
        cd  /tmp
        curl -s -LO https://storage.googleapis.com/golang/go#{GO_VERSION}.linux-amd64.tar.gz
        tar xzf go#{GO_VERSION}.linux-amd64.tar.gz
        mv go /usr/local/
        echo 'export GOPATH=$HOME/go' >> /home/vagrant/.profile
        echo 'export PATH=$HOME/go/bin:/usr/local/go/bin:$PATH' >> /home/vagrant/.profile
    fi
 SHELL
end
