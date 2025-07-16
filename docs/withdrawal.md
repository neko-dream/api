# User Withdrawal Functionality

## Overview

This document describes the user withdrawal functionality implemented to allow users to properly withdraw from the system while maintaining data integrity and preventing system abuse.

## Features

### Data Handling on Withdrawal

**Data that is deleted/anonymized:**
- User's display ID (replaced with "deleted_user")
- User's display name (replaced with "削除されたユーザー")
- User's icon URL (set to null)
- User's authentication subject (withdrawal date is recorded)

**Data that is preserved:**
- Statistical/demographic information (age, gender, etc.) for analysis purposes
- User opinions and session participation (but anonymized)
- Login history for audit purposes

### Re-registration Restrictions

- Users cannot re-register for 30 days after withdrawal
- This prevents abuse of the voting system by repeatedly withdrawing and re-registering
- The restriction is enforced by checking the withdrawal date during authentication

### Technical Implementation

#### Database Changes
- Added `withdrawal_date` column to `user_auths` table
- Added index for efficient querying of withdrawn users

#### API Endpoints
- Updated `/auth/dev/detach` endpoint to perform proper withdrawal instead of temporary detachment
- The endpoint now:
  1. Anonymizes user data
  2. Records withdrawal date
  3. Deactivates all active sessions
  4. Revokes current authentication

#### Domain Models
- `User` domain model includes withdrawal functionality
- `UserAuth` domain model tracks withdrawal status and re-registration eligibility
- Proper encapsulation of withdrawal business logic

## Usage

### Withdrawing a User
```
DELETE /auth/dev/detach
```

The authenticated user will be withdrawn from the system. Their data will be anonymized and they will be logged out.

### Re-registration Process
When a user attempts to log in after withdrawal:
1. The system checks if the withdrawal date exists
2. If less than 30 days have passed, login is denied
3. If 30+ days have passed, normal login proceeds

## Database Migration

The withdrawal functionality requires running migration `000039_add_withdrawal_date_to_user_auths.up.sql` which adds the necessary database field.

## Security Considerations

- Withdrawal is irreversible - user data cannot be recovered
- Statistical data is preserved for analysis but is fully anonymized
- Re-registration restrictions prevent voting system abuse
- All active sessions are terminated on withdrawal