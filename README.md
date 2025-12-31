# DevPool : The Bot for Assigning n' Slicing Contributors (Season of Code)

The bot with a mouth at Summer of Code. DevPool does not listen 
to any of your commands by itself and gets all its commands through a Redis 
stream that it has subscribed to. All your commands are captured by 
[Alfred - The Webhook Butler](https://github.com/Infinite-Sum-Games/alfred). 

Read more about him at his home repository!

## Command Workflows

### 1. `/assign` - Contributor ONLY Command
Use this command to assign an issue to yourself. Your must be previously 
registered as a participant into the Summer of Code portal for this command 
to work. Additionally, this command only works on issues that are accepted for 
the program.

### 2. `/unassign` - Contributor ONLY Command
Use this command to un-assign yourself from an issue. Your must be previously 
registered as a participant into the Summer of Code portal for this command 
to work. Additionally, this command only works on issues that are accepted for 
the program.

### 3. `/bounty <amount> @github-username` - Maintainer ONLY Command
Use this command to give out a bounty to a contributor. Bear note, the bounty 
will be given only if the contributor is registered to the program. It would 
otherwise fail and the bot would return with an unhappy message :(

### 4. `/penalty <amount> @github-username` - Maintainer ONLY Command
Use this command to punish contributors for misbehavior. Working pattern is 
same as the previously mentioned bounty command.

### 5. `/bug @github-username` - Maintainer ONLY Command
Use this command in an "issue" only to accept a bug report. This allows a 
contributor to be eligible for the "bug-hunting" related achievement badges.

### 6. `/impact @github-username` - Maintainer ONLY Command
Use this command to mark a pull request as a "high-impact" contribution so that 
the contributor is eligible for the corresponding achievement badge.

### 7. `/doc @github-username` - Maintainer ONLY Command
Use this command to mark a pull request as a "documentation" contribution. This
allows a contributor to be eligible for documentation-related achievement 
badges.

### 8. `/test @github-username` - Maintainer ONLY Command
Use this command to mark a pull request as a "testing" contribution. This allows
a contributor to be eligible for testing-related achievement badges.

### 9. `/help @github-username` - Maintainer ONLY Command
Use this command to mark a contributor as a helper. Can be used in both issues 
as well as pull-requests.

## Development Instructions
1. Clone the repository
```bash
git clone https://github.com/Infinite-Sum-Games/devpool
```
2. Populate the `.env.example` file with appropriate credentials and 
rename it to `.env`.
3. Setup Redis / Valkey along with Redis-Insight. This setup is automatically 
handled for you if you are running `docker-compose` from 
[Infinite-Sum-Games/pulse](https://github.com/Infinite-Sum-Games/pulse).
4. Load the sample data to test within Redis-Insight to check. Make sure to 
clean it out.

## Authors
This bot is developed and tested by [Ritesh Koushik](https://github.com/IAmRiteshKoushik).
