type Platform = "twitch" | "youtube";

type User = {
  id: string;
  name: string;
  broadcaster: boolean;
  moderator: boolean;
  twitch_vip: boolean;
  youtube_member: boolean;
};

export type AdminWSMessage = {
  id: string;
  body: string;
  platform: Platform;
  sender: User;
  received_at: Date;
  published_at: Date;
};