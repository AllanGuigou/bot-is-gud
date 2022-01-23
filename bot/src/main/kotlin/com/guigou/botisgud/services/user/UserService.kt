package com.guigou.botisgud.services.user

import dev.kord.common.entity.Snowflake
import com.guigou.botisgud.models.User

interface UserService {
    suspend fun add(user: User)
    suspend fun get(userId: Snowflake): User?
}
