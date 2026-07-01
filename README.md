Gator is a simple RSS agregator. It's a guided project for boot.dev.

It uses PostgreSQL and Go. So you need to install Postgres and Go.

To install Gator navigate to a root of a program and run:
$go intall
or use path to a program:
&go intall path/to/root/Gator

For program to work you to setup config file. Go to your user's home directoty like "/home/user" and make ".gatorconfig.json" file.
Add:
{"db_url":"postgres://postgres:postgres@localhost:5432/gator?sslmode=disable","current_user_name":""}
to a ".gatorconfig.json".

To use Gator write:
gator <name of command>

List of commands:
"register" <username> - register a new user and set it logged in.
"login" <username> - login as already existing user.
"addfeed" <name> <url> - add new feed to a datebase.
"agg" - start to fetch posts from feeds. Use "ctrl+c" to stop it.
"browse" <limit> - print <limit> posts for the current user.