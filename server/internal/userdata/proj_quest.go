package userdata

import (
	"sort"

	"lunar-tear/server/internal/store"
	"lunar-tear/server/internal/utils"
)

func sortedQuestRecords(user store.UserState) []map[string]any {
	ids := make([]int, 0, len(user.Quests))
	for id := range user.Quests {
		ids = append(ids, int(id))
	}
	sort.Ints(ids)
	records := make([]map[string]any, 0, len(ids))
	for _, id := range ids {
		row := user.Quests[int32(id)]
		records = append(records, map[string]any{
			"userId":              user.UserId,
			"questId":             row.QuestId,
			"questStateType":      row.QuestStateType,
			"isBattleOnly":        row.IsBattleOnly,
			"latestStartDatetime": row.LatestStartDatetime,
			"clearCount":          row.ClearCount,
			"dailyClearCount":     row.DailyClearCount,
			"lastClearDatetime":   row.LastClearDatetime,
			"shortestClearFrames": row.ShortestClearFrames,
			"latestVersion":       row.LatestVersion,
		})
	}
	return records
}

func sortedQuestMissionRecords(user store.UserState) []map[string]any {
	keys := make([]store.QuestMissionKey, 0, len(user.QuestMissions))
	for key := range user.QuestMissions {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].QuestId != keys[j].QuestId {
			return keys[i].QuestId < keys[j].QuestId
		}
		return keys[i].QuestMissionId < keys[j].QuestMissionId
	})
	records := make([]map[string]any, 0, len(keys))
	for _, key := range keys {
		row := user.QuestMissions[key]
		records = append(records, map[string]any{
			"userId":              user.UserId,
			"questId":             row.QuestId,
			"questMissionId":      row.QuestMissionId,
			"progressValue":       row.ProgressValue,
			"isClear":             row.IsClear,
			"latestClearDatetime": row.LatestClearDatetime,
			"latestVersion":       row.LatestVersion,
		})
	}
	return records
}

func init() {
	register("IUserQuest", func(user store.UserState) string {
		s, _ := utils.EncodeJSONMaps(sortedQuestRecords(user)...)
		return s
	})
	register("IUserQuestMission", func(user store.UserState) string {
		s, _ := utils.EncodeJSONMaps(sortedQuestMissionRecords(user)...)
		return s
	})
	register("IUserMainQuestFlowStatus", func(user store.UserState) string {
		s, _ := utils.EncodeJSONMaps(map[string]any{
			"userId":               user.UserId,
			"currentQuestFlowType": user.MainQuest.CurrentQuestFlowType,
			"latestVersion":        user.MainQuest.LatestVersion,
		})
		return s
	})
	register("IUserMainQuestMainFlowStatus", func(user store.UserState) string {
		s, _ := utils.EncodeJSONMaps(map[string]any{
			"userId":                  user.UserId,
			"currentMainQuestRouteId": user.MainQuest.CurrentMainQuestRouteId,
			"currentQuestSceneId":     user.MainQuest.CurrentQuestSceneId,
			"headQuestSceneId":        user.MainQuest.HeadQuestSceneId,
			"isReachedLastQuestScene": user.MainQuest.IsReachedLastQuestScene,
			"latestVersion":           user.MainQuest.LatestVersion,
		})
		return s
	})
	register("IUserMainQuestProgressStatus", func(user store.UserState) string {
		s, _ := utils.EncodeJSONMaps(map[string]any{
			"userId":               user.UserId,
			"currentQuestSceneId":  user.MainQuest.ProgressQuestSceneId,
			"headQuestSceneId":     user.MainQuest.ProgressHeadQuestSceneId,
			"currentQuestFlowType": user.MainQuest.ProgressQuestFlowType,
			"latestVersion":        user.MainQuest.LatestVersion,
		})
		return s
	})
	register("IUserMainQuestSeasonRoute", func(user store.UserState) string {
		if len(user.MainQuestSeasonRoutes) == 0 {
			// Fallback to current (season, route) for legacy saves with no history.
			s, _ := utils.EncodeJSONMaps(map[string]any{
				"userId":            user.UserId,
				"mainQuestSeasonId": user.MainQuest.MainQuestSeasonId,
				"mainQuestRouteId":  user.MainQuest.CurrentMainQuestRouteId,
				"latestVersion":     user.MainQuest.LatestVersion,
			})
			return s
		}
		keys := make([]store.SeasonRouteKey, 0, len(user.MainQuestSeasonRoutes))
		for k := range user.MainQuestSeasonRoutes {
			keys = append(keys, k)
		}
		sort.Slice(keys, func(i, j int) bool {
			if keys[i].MainQuestSeasonId != keys[j].MainQuestSeasonId {
				return keys[i].MainQuestSeasonId < keys[j].MainQuestSeasonId
			}
			return keys[i].MainQuestRouteId < keys[j].MainQuestRouteId
		})
		records := make([]map[string]any, 0, len(keys))
		for _, k := range keys {
			e := user.MainQuestSeasonRoutes[k]
			records = append(records, map[string]any{
				"userId":            user.UserId,
				"mainQuestSeasonId": e.MainQuestSeasonId,
				"mainQuestRouteId":  e.MainQuestRouteId,
				"latestVersion":     e.LatestVersion,
			})
		}
		s, _ := utils.EncodeJSONMaps(records...)
		return s
	})
	register("IUserEventQuestProgressStatus", func(user store.UserState) string {
		s, _ := utils.EncodeJSONMaps(map[string]any{
			"userId":                     user.UserId,
			"currentEventQuestChapterId": user.EventQuest.CurrentEventQuestChapterId,
			"currentQuestId":             user.EventQuest.CurrentQuestId,
			"currentQuestSceneId":        user.EventQuest.CurrentQuestSceneId,
			"headQuestSceneId":           user.EventQuest.HeadQuestSceneId,
			"latestVersion":              user.EventQuest.LatestVersion,
		})
		return s
	})
	register("IUserExtraQuestProgressStatus", func(user store.UserState) string {
		s, _ := utils.EncodeJSONMaps(map[string]any{
			"userId":              user.UserId,
			"currentQuestId":      user.ExtraQuest.CurrentQuestId,
			"currentQuestSceneId": user.ExtraQuest.CurrentQuestSceneId,
			"headQuestSceneId":    user.ExtraQuest.HeadQuestSceneId,
			"latestVersion":       user.ExtraQuest.LatestVersion,
		})
		return s
	})
	register("IUserMainQuestReplayFlowStatus", func(user store.UserState) string {
		if user.MainQuest.ReplayFlowCurrentQuestSceneId == 0 && user.MainQuest.ReplayFlowHeadQuestSceneId == 0 {
			return "[]"
		}
		s, _ := utils.EncodeJSONMaps(map[string]any{
			"userId":                  user.UserId,
			"currentHeadQuestSceneId": user.MainQuest.ReplayFlowHeadQuestSceneId,
			"currentQuestSceneId":     user.MainQuest.ReplayFlowCurrentQuestSceneId,
			"latestVersion":           user.MainQuest.LatestVersion,
		})
		return s
	})
	register("IUserSideStoryQuestSceneProgressStatus", func(user store.UserState) string {
		s, _ := utils.EncodeJSONMaps(map[string]any{
			"userId":                       user.UserId,
			"currentSideStoryQuestId":      user.SideStoryActiveProgress.CurrentSideStoryQuestId,
			"currentSideStoryQuestSceneId": user.SideStoryActiveProgress.CurrentSideStoryQuestSceneId,
			"latestVersion":                user.SideStoryActiveProgress.LatestVersion,
		})
		return s
	})
	register("IUserSideStoryQuest", func(user store.UserState) string {
		if len(user.SideStoryQuests) == 0 {
			return "[]"
		}
		ids := make([]int, 0, len(user.SideStoryQuests))
		for id := range user.SideStoryQuests {
			ids = append(ids, int(id))
		}
		sort.Ints(ids)
		records := make([]map[string]any, 0, len(ids))
		for _, id := range ids {
			progress := user.SideStoryQuests[int32(id)]
			records = append(records, map[string]any{
				"userId":                    user.UserId,
				"sideStoryQuestId":          int32(id),
				"headSideStoryQuestSceneId": progress.HeadSideStoryQuestSceneId,
				"sideStoryQuestStateType":   progress.SideStoryQuestStateType,
				"latestVersion":             progress.LatestVersion,
			})
		}
		s, _ := utils.EncodeJSONMaps(records...)
		return s
	})
	register("IUserQuestLimitContentStatus", func(user store.UserState) string {
		if len(user.QuestLimitContentStatus) == 0 {
			return "[]"
		}
		ids := make([]int, 0, len(user.QuestLimitContentStatus))
		for id := range user.QuestLimitContentStatus {
			ids = append(ids, int(id))
		}
		sort.Ints(ids)
		records := make([]map[string]any, 0, len(ids))
		for _, id := range ids {
			st := user.QuestLimitContentStatus[int32(id)]
			records = append(records, map[string]any{
				"userId":                      user.UserId,
				"questId":                     int32(id),
				"limitContentQuestStatusType": st.LimitContentQuestStatusType,
				"eventQuestChapterId":         st.EventQuestChapterId,
				"latestVersion":               st.LatestVersion,
			})
		}
		s, _ := utils.EncodeJSONMaps(records...)
		return s
	})
	register("IUserEventQuestTowerAccumulationReward", func(user store.UserState) string {
		if len(user.TowerAccumulationRewards) == 0 {
			return "[]"
		}
		ids := make([]int, 0, len(user.TowerAccumulationRewards))
		for id := range user.TowerAccumulationRewards {
			ids = append(ids, int(id))
		}
		sort.Ints(ids)
		records := make([]map[string]any, 0, len(ids))
		for _, id := range ids {
			st := user.TowerAccumulationRewards[int32(id)]
			records = append(records, map[string]any{
				"userId":              user.UserId,
				"eventQuestChapterId": st.EventQuestChapterId,
				"latestRewardReceiveQuestMissionClearCount": st.LatestRewardReceiveQuestMissionClearCount,
				"latestVersion": st.LatestVersion,
			})
		}
		s, _ := utils.EncodeJSONMaps(records...)
		return s
	})
	registerStatic(
		"IUserEventQuestDailyGroupCompleteReward",
		"IUserQuestReplayFlowRewardGroup",
		"IUserQuestAutoOrbit",
		"IUserQuestSceneChoice",
		"IUserQuestSceneChoiceHistory",
	)
}
