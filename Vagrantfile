# -*- mode: ruby -*-
# vi: set ft=ruby :

MOUNT_POINT = "/vagrant"

Vagrant::Config.run do |config|
  # Every Vagrant virtual environment requires a box to build off of.
  config.vm.box = "centos/7"
  config.vm.host_name = "go-deploy-register"

  #config.vm.forward_port 80, 80 # Nginx
  config.vm.network "forwarded_port", host: 80, guest: 80, host_ip: '0.0.0.0'

  # Configure a private network required by nfs folder share
  config.vm.network :hostonly, "33.33.33.190"

  # Root shared folder using nfs for better performance
  config.vm.share_folder("v-root", MOUNT_POINT, ".", :nfs => true)
end
