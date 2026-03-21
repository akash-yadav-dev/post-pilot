import type {
  ActivityItem,
  AnalyticsPoint,
  CreatePostInput,
  DashboardMetrics,
  MediaAsset,
  NotificationItem,
  PlatformSummary,
  PostEntity,
  PostStatus,
  Slot,
  TeamMember,
} from "@/lib/social/types";

const now = new Date();

const db: {
  posts: PostEntity[];
  activities: ActivityItem[];
  notifications: NotificationItem[];
  team: TeamMember[];
  assets: MediaAsset[];
  slots: Slot[];
  analytics: AnalyticsPoint[];
  platforms: PlatformSummary[];
} = {
  posts: [
    {
      id: "post_1",
      title: "Launch teaser",
      content: "New release incoming. Stay tuned. #product",
      mediaUrls: ["https://images.unsplash.com/photo-1520607162513-77705c0f0d4a?w=800"],
      platforms: ["twitter", "linkedin"],
      scheduledAt: new Date(now.getTime() + 1000 * 60 * 60 * 6).toISOString(),
      status: "scheduled",
      likes: 46,
      comments: 7,
      shares: 11,
      createdAt: new Date(now.getTime() - 1000 * 60 * 60 * 10).toISOString(),
      updatedAt: new Date(now.getTime() - 1000 * 60 * 60 * 3).toISOString(),
    },
    {
      id: "post_2",
      title: "Case study",
      content: "How social ops teams cut planning time by 41%.",
      mediaUrls: [],
      platforms: ["linkedin"],
      status: "draft",
      likes: 0,
      comments: 0,
      shares: 0,
      createdAt: new Date(now.getTime() - 1000 * 60 * 60 * 40).toISOString(),
      updatedAt: new Date(now.getTime() - 1000 * 60 * 60 * 6).toISOString(),
    },
    {
      id: "post_3",
      title: "Friday roundup",
      content: "Top growth experiments from this week.",
      mediaUrls: ["https://images.unsplash.com/photo-1551281044-8b1a4f3f5f8b?w=800"],
      platforms: ["facebook", "instagram"],
      status: "published",
      scheduledAt: new Date(now.getTime() - 1000 * 60 * 60 * 16).toISOString(),
      likes: 312,
      comments: 28,
      shares: 44,
      createdAt: new Date(now.getTime() - 1000 * 60 * 60 * 48).toISOString(),
      updatedAt: new Date(now.getTime() - 1000 * 60 * 60 * 15).toISOString(),
    },
    {
      id: "post_4",
      title: "A/B test recap",
      content: "Results from headline test B.",
      mediaUrls: [],
      platforms: ["twitter"],
      status: "failed",
      scheduledAt: new Date(now.getTime() - 1000 * 60 * 60 * 2).toISOString(),
      likes: 0,
      comments: 0,
      shares: 0,
      createdAt: new Date(now.getTime() - 1000 * 60 * 60 * 22).toISOString(),
      updatedAt: new Date(now.getTime() - 1000 * 60 * 60).toISOString(),
    },
  ],
  activities: [
    { id: "a1", action: "Scheduled post for LinkedIn", actor: "Akash", time: "4 min ago" },
    { id: "a2", action: "Approved post in review queue", actor: "Sonia", time: "11 min ago" },
    { id: "a3", action: "Connected Instagram account", actor: "Akash", time: "1 hr ago" },
  ],
  notifications: [
    { id: "n1", title: "Post published", description: "Friday roundup is now live.", kind: "success", read: false, createdAt: "2m ago" },
    { id: "n2", title: "Approval requested", description: "Case study draft needs review.", kind: "info", read: false, createdAt: "9m ago" },
    { id: "n3", title: "Publish failed", description: "A/B test recap failed on X due to token expiry.", kind: "warning", read: true, createdAt: "1h ago" },
  ],
  team: [
    { id: "t1", name: "Akash Yadav", role: "admin", email: "akash@postpilot.app" },
    { id: "t2", name: "Sonia Price", role: "editor", email: "sonia@postpilot.app" },
    { id: "t3", name: "Mark Lee", role: "viewer", email: "mark@postpilot.app" },
  ],
  assets: [
    { id: "m1", name: "campaign-cover.jpg", url: "https://images.unsplash.com/photo-1519389950473-47ba0277781c?w=500", type: "image", tags: ["launch", "hero"], uploadedAt: "2d ago" },
    { id: "m2", name: "teaser-video.mp4", url: "https://images.unsplash.com/photo-1498050108023-c5249f4df085?w=500", type: "video", tags: ["video", "teaser"], uploadedAt: "5d ago" },
  ],
  slots: [
    { id: "s1", label: "Morning", time: "09:00", timezone: "Asia/Kolkata" },
    { id: "s2", label: "Noon", time: "13:00", timezone: "Asia/Kolkata" },
    { id: "s3", label: "Evening", time: "18:30", timezone: "Asia/Kolkata" },
  ],
  analytics: [
    { label: "Mon", engagementRate: 3.8, reach: 8120, ctr: 2.1 },
    { label: "Tue", engagementRate: 4.2, reach: 9230, ctr: 2.4 },
    { label: "Wed", engagementRate: 4.5, reach: 10110, ctr: 2.8 },
    { label: "Thu", engagementRate: 5.0, reach: 11020, ctr: 3.2 },
    { label: "Fri", engagementRate: 5.4, reach: 12900, ctr: 3.7 },
  ],
  platforms: [
    { platform: "twitter", connected: true, posts: 36, engagementRate: 4.4, health: "healthy" },
    { platform: "linkedin", connected: true, posts: 22, engagementRate: 5.6, health: "healthy" },
    { platform: "facebook", connected: true, posts: 18, engagementRate: 3.1, health: "warning" },
    { platform: "instagram", connected: true, posts: 27, engagementRate: 6.3, health: "healthy" },
    { platform: "pinterest", connected: false, posts: 0, engagementRate: 0, health: "expired" },
  ],
};

function delay(ms = 220) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

export const SocialDashboardService = {
  async getMetrics(): Promise<DashboardMetrics> {
    await delay();
    const scheduledCount = db.posts.filter((p) => p.status === "scheduled" || p.status === "queued").length;
    const publishedCount = db.posts.filter((p) => p.status === "published").length;
    const engagement = db.posts.reduce((sum, p) => sum + p.likes + p.comments + p.shares, 0);
    const failedCount = db.posts.filter((p) => p.status === "failed").length;
    return { scheduledCount, publishedCount, engagement, failedCount };
  },

  async getPosts(params?: { status?: PostStatus; platform?: string; search?: string }): Promise<PostEntity[]> {
    await delay();
    let posts = [...db.posts];
    if (params?.status) {
      posts = posts.filter((p) => p.status === params.status);
    }
    if (params?.platform) {
      posts = posts.filter((p) => p.platforms.includes(params.platform as never));
    }
    if (params?.search) {
      const q = params.search.toLowerCase();
      posts = posts.filter((p) => p.title.toLowerCase().includes(q) || p.content.toLowerCase().includes(q));
    }
    return posts.sort((a, b) => b.updatedAt.localeCompare(a.updatedAt));
  },

  async getUpcomingPosts(): Promise<PostEntity[]> {
    await delay();
    return db.posts
      .filter((p) => p.status === "scheduled" || p.status === "queued")
      .sort((a, b) => (a.scheduledAt ?? "").localeCompare(b.scheduledAt ?? ""));
  },

  async getPlatformSummary(): Promise<PlatformSummary[]> {
    await delay();
    return [...db.platforms];
  },

  async getActivityFeed(): Promise<ActivityItem[]> {
    await delay();
    return [...db.activities];
  },

  async getNotifications(): Promise<NotificationItem[]> {
    await delay();
    return [...db.notifications];
  },

  async getTeamMembers(): Promise<TeamMember[]> {
    await delay();
    return [...db.team];
  },

  async getAssets(): Promise<MediaAsset[]> {
    await delay();
    return [...db.assets];
  },

  async getQueueSlots(): Promise<Slot[]> {
    await delay();
    return [...db.slots];
  },

  async getAnalytics(): Promise<AnalyticsPoint[]> {
    await delay();
    return [...db.analytics];
  },

  async createPost(input: CreatePostInput): Promise<PostEntity> {
    await delay();
    const post: PostEntity = {
      id: `post_${Math.random().toString(36).slice(2, 9)}`,
      title: input.title,
      content: input.content,
      mediaUrls: input.mediaUrls,
      platforms: input.platforms,
      scheduledAt: input.scheduledAt,
      status: input.status,
      likes: 0,
      comments: 0,
      shares: 0,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    };
    db.posts.unshift(post);
    db.activities.unshift({ id: `a_${Date.now()}`, actor: "You", action: `Created ${post.status} post`, time: "just now" });
    db.notifications.unshift({
      id: `n_${Date.now()}`,
      title: "Post saved",
      description: `\"${post.title}\" was saved as ${post.status}.`,
      kind: "success",
      read: false,
      createdAt: "now",
    });
    return post;
  },

  async updatePostStatus(id: string, status: PostStatus): Promise<PostEntity> {
    await delay();
    const target = db.posts.find((p) => p.id === id);
    if (!target) {
      throw new Error("Post not found");
    }
    target.status = status;
    target.updatedAt = new Date().toISOString();
    return target;
  },

  async deletePost(id: string): Promise<void> {
    await delay();
    db.posts = db.posts.filter((p) => p.id !== id);
  },

  async duplicatePost(id: string): Promise<PostEntity> {
    await delay();
    const source = db.posts.find((p) => p.id === id);
    if (!source) {
      throw new Error("Post not found");
    }
    const duplicate: PostEntity = {
      ...source,
      id: `post_${Math.random().toString(36).slice(2, 9)}`,
      title: `${source.title} (Copy)`,
      status: "draft",
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    };
    db.posts.unshift(duplicate);
    return duplicate;
  },

  async markNotificationRead(id: string): Promise<void> {
    await delay(120);
    const n = db.notifications.find((item) => item.id === id);
    if (n) {
      n.read = true;
    }
  },
};
