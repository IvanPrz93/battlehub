# BattleHub

BattleHub is a backend service for a fictional competitive PvP game.

## Features

- User Authentication
- Player Profiles
- Match Creation
- Match Results
- Match History
- ELO Ranking System
- Leaderboard

## Authentication

Users can:
- Register
- Login

Authentication uses JWT.

## Player Profiles

Players can:
- View their own profile
- View another player's profile by username

Profiles include:
- Wins
- Losses
- Rating

## Match System

Matches can be:
- Created
- Completed
- Stored in match history

When a match is completed:
- The winner is recorded
- Wins and losses are updated
- ELO ratings are recalculated

Match completion is handled using a database transaction to ensure consistency.

## Match History

Every match result is stored and users can view:
- Their own match history
- Another player's match history

## Leaderboard

BattleHub uses an ELO rating system. Wins and losses change rating based on opponent rating.

Players are ranked based on rating and can be viewed through the leaderboard endpoint.