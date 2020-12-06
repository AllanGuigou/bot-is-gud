package com.guigou.botisgud.services.user

import com.gitlab.kordlib.common.entity.Snowflake
import com.guigou.botisgud.extensions.logger
import com.guigou.botisgud.models.User
import kotlinx.coroutines.Dispatchers
import org.jetbrains.exposed.sql.transactions.experimental.newSuspendedTransaction

class UserServiceImpl : UserService {
    companion object {
        val logger = logger()
    }

    override suspend fun add(user: User) {
        newSuspendedTransaction(Dispatchers.Default) {
            UserEntity.new {
                userId = user.userId
                zone = user.zone
            }
        }
    }

    override suspend fun get(userId: Snowflake): User? {
        return newSuspendedTransaction(Dispatchers.Default) {
            val user = UserEntity.find { Users.userId eq userId.longValue }.firstOrNull()
                ?: return@newSuspendedTransaction null

            return@newSuspendedTransaction User(user.userId, user.zone)
        }
    }
}
