"use client";

import {
  keepPreviousData,
  useMutation,
  useQuery,
  useQueryClient,
  type UseQueryResult,
} from "@tanstack/react-query";
import { SocialDashboardService } from "@/services/social-dashboard-service";
import type {
  CreatePostInput,
  NotificationItem,
  PostEntity,
  PostStatus,
} from "@/lib/social/types";

export const socialQueryKeys = {
  metrics: ["social", "metrics"] as const,
  upcoming: ["social", "upcoming"] as const,
  activity: ["social", "activity"] as const,
  platformSummary: ["social", "platform-summary"] as const,
  analytics: ["social", "analytics"] as const,
  notifications: ["social", "notifications"] as const,
  team: ["social", "team"] as const,
  assets: ["social", "assets"] as const,
  queueSlots: ["social", "queue-slots"] as const,
  posts: (filters: { status?: PostStatus; platform?: string; search?: string }) =>
    ["social", "posts", filters] as const,
};

export function useDashboardMetrics() {
  return useQuery({ queryKey: socialQueryKeys.metrics, queryFn: SocialDashboardService.getMetrics });
}

export function useUpcomingPosts() {
  return useQuery({ queryKey: socialQueryKeys.upcoming, queryFn: SocialDashboardService.getUpcomingPosts });
}

export function usePlatformSummary() {
  return useQuery({ queryKey: socialQueryKeys.platformSummary, queryFn: SocialDashboardService.getPlatformSummary });
}

export function useActivityFeed() {
  return useQuery({ queryKey: socialQueryKeys.activity, queryFn: SocialDashboardService.getActivityFeed });
}

export function usePosts(filters: { status?: PostStatus; platform?: string; search?: string }): UseQueryResult<PostEntity[]> {
  return useQuery({
    queryKey: socialQueryKeys.posts(filters),
    queryFn: () => SocialDashboardService.getPosts(filters),
    placeholderData: keepPreviousData,
  });
}

export function useAnalytics() {
  return useQuery({ queryKey: socialQueryKeys.analytics, queryFn: SocialDashboardService.getAnalytics });
}

export function useNotifications() {
  return useQuery({ queryKey: socialQueryKeys.notifications, queryFn: SocialDashboardService.getNotifications });
}

export function useTeamMembers() {
  return useQuery({ queryKey: socialQueryKeys.team, queryFn: SocialDashboardService.getTeamMembers });
}

export function useMediaAssets() {
  return useQuery({ queryKey: socialQueryKeys.assets, queryFn: SocialDashboardService.getAssets });
}

export function useQueueSlots() {
  return useQuery({ queryKey: socialQueryKeys.queueSlots, queryFn: SocialDashboardService.getQueueSlots });
}

export function useCreatePostMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (input: CreatePostInput) => SocialDashboardService.createPost(input),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["social"] });
    },
  });
}

export function useDeletePostMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => SocialDashboardService.deletePost(id),
    onMutate: async (id) => {
      await queryClient.cancelQueries({ queryKey: ["social", "posts"] });
      const snapshots = queryClient.getQueriesData<PostEntity[]>({ queryKey: ["social", "posts"] });
      snapshots.forEach(([key, posts]) => {
        if (!posts) return;
        queryClient.setQueryData<PostEntity[]>(key, posts.filter((post) => post.id !== id));
      });
      return { snapshots };
    },
    onError: (_err, _id, context) => {
      context?.snapshots?.forEach(([key, data]) => {
        queryClient.setQueryData(key, data);
      });
    },
    onSettled: () => {
      void queryClient.invalidateQueries({ queryKey: ["social"] });
    },
  });
}

export function useDuplicatePostMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => SocialDashboardService.duplicatePost(id),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["social"] });
    },
  });
}

export function useUpdatePostStatusMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, status }: { id: string; status: PostStatus }) =>
      SocialDashboardService.updatePostStatus(id, status),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["social"] });
    },
  });
}

export function useMarkNotificationReadMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => SocialDashboardService.markNotificationRead(id),
    onMutate: async (id: string) => {
      await queryClient.cancelQueries({ queryKey: socialQueryKeys.notifications });
      const previous = queryClient.getQueryData<NotificationItem[]>(socialQueryKeys.notifications);
      if (previous) {
        queryClient.setQueryData<NotificationItem[]>(
          socialQueryKeys.notifications,
          previous.map((item) => (item.id === id ? { ...item, read: true } : item))
        );
      }
      return { previous };
    },
    onError: (_err, _id, context) => {
      if (context?.previous) {
        queryClient.setQueryData(socialQueryKeys.notifications, context.previous);
      }
    },
    onSettled: () => {
      void queryClient.invalidateQueries({ queryKey: socialQueryKeys.notifications });
    },
  });
}
