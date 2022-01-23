package com.guigou.botisgud.services.nickname

import dev.kord.common.entity.Snowflake

interface NicknameService {
    suspend fun getUsersByGuild(): Map<Snowflake, List<Snowflake>>
}
