export type SocialPlatform = "twitter" | "linkedin" | "facebook" | "instagram" | "pinterest";

export type PostStatus = "draft" | "scheduled" | "published" | "failed" | "queued" | "review" | "approved";

export type PostEntity = {
  id: string;
  title: string;
  content: string;
  mediaUrls: string[];
  platforms: SocialPlatform[];
  scheduledAt?: string;
  status: PostStatus;
  likes: number;
  comments: number;
  shares: number;
  createdAt: string;
  updatedAt: string;
};

export type DashboardMetrics = {
  scheduledCount: number;
  publishedCount: number;
  engagement: number;
  failedCount: number;
};

export type PlatformSummary = {
  platform: SocialPlatform;
  connected: boolean;
  posts: number;
  engagementRate: number;
  health: "healthy" | "warning" | "expired";
};

export type ActivityItem = {
  id: string;
  action: string;
  actor: string;
  time: string;
};

export type NotificationItem = {
  id: string;
  title: string;
  description: string;
  kind: "success" | "warning" | "info";
  read: boolean;
  createdAt: string;
};

export type TeamMember = {
  id: string;
  name: string;
  role: "admin" | "editor" | "viewer";
  email: string;
};

export type MediaAsset = {
  id: string;
  name: string;
  url: string;
  type: "image" | "video";
  tags: string[];
  uploadedAt: string;
};

export type Slot = {
  id: string;
  label: string;
  time: string;
  timezone: string;
};

export type AnalyticsPoint = {
  label: string;
  engagementRate: number;
  reach: number;
  ctr: number;
};

export type CreatePostInput = {
  title: string;
  content: string;
  mediaUrls: string[];
  platforms: SocialPlatform[];
  scheduledAt?: string;
  status: PostStatus;
};
