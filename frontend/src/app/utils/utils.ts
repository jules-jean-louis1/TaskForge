export enum USER_ROLE {
  ADMIN = 'admin',
  TECH = 'tech',
  STANDARD = 'standard',
}

export const USER_AVAILABLES_ROLES = [USER_ROLE.ADMIN, USER_ROLE.TECH, USER_ROLE.STANDARD] as const;

export type UserRole = (typeof USER_AVAILABLES_ROLES)[number];
