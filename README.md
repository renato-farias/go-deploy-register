# go-deploy-register

## Preparing the environment

```
vagrant up
```

We need Go compiler for running and building our API. So we hope Go, MySQL server are installed in your environment.

```
setup/setup.sh
```

Now set the $GOPATH variable creating a new folder in your home. $GOPATH is where your external packages and libraries are stored. My ~/.bashrc looks like this.

```
export GOPATH=/vagrant
```

then, reload your environments.

```
source ~/.bashrc
```

Now we are ready with setup. Let us install dependencies.

```
go get "github.com/go-sql-driver/mysql"
go get "github.com/gin-gonic/gin"
go get "github.com/jinzhu/configor"
```

sudo service mariadb start

sudo mysql_secure_installation

[vagrant@go-deploy-register app]$ sudo mysql_secure_installation

NOTE: RUNNING ALL PARTS OF THIS SCRIPT IS RECOMMENDED FOR ALL MariaDB
      SERVERS IN PRODUCTION USE!  PLEASE READ EACH STEP CAREFULLY!

In order to log into MariaDB to secure it, we'll need the current
password for the root user.  If you've just installed MariaDB, and
you haven't set the root password yet, the password will be blank,
so you should just press enter here.

Enter current password for root (enter for none):
OK, successfully used password, moving on...

Setting the root password ensures that nobody can log into the MariaDB
root user without the proper authorisation.

Set root password? [Y/n] Y
New password:
Re-enter new password:
Password updated successfully!
Reloading privilege tables..
 ... Success!


By default, a MariaDB installation has an anonymous user, allowing anyone
to log into MariaDB without having to have a user account created for
them.  This is intended only for testing, and to make the installation
go a bit smoother.  You should remove them before moving into a
production environment.

Remove anonymous users? [Y/n] Y
 ... Success!

Normally, root should only be allowed to connect from 'localhost'.  This
ensures that someone cannot guess at the root password from the network.

Disallow root login remotely? [Y/n] Y
 ... Success!

By default, MariaDB comes with a database named 'test' that anyone can
access.  This is also intended only for testing, and should be removed
before moving into a production environment.

Remove test database and access to it? [Y/n] Y
 - Dropping test database...
 ... Success!
 - Removing privileges on test database...
 ... Success!

Reloading the privilege tables will ensure that all changes made so far
will take effect immediately.

Reload privilege tables now? [Y/n] Y
 ... Success!

Cleaning up...

All done!  If you've completed all of the above steps, your MariaDB
installation should now be secure.

Thanks for using MariaDB!

[vagrant@go-deploy-register app]$ mysql -u root -p
Enter password:


# CREATE USER 'deploy_register'@'%' IDENTIFIED BY 'deploy_register';
# CREATE DATABASE deploy_register;
# GRANT ALL PRIVILEGES ON *.* TO 'deploy_register'@'%' WITH GRANT OPTION;

http://33.33.33.190:3000/application/1

http://33.33.33.190:3000/applications

curl http://33.33.33.190:3000/application -d 'application_name=main_application' --user 'deploy:secretpassword'

curl http://33.33.33.190:3000/deployment -d 'version=teste&application_id=1' --user 'deploy:secretpassword'
