package com.guigou.botisgud.services.nickname

import com.gitlab.kordlib.common.entity.Snowflake

interface NicknameService {
    suspend fun getUsersByGuild(): Map<Snowflake, List<Snowflake>>
}
