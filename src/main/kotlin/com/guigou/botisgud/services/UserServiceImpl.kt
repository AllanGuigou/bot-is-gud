package com.guigou.botisgud.services

import com.gitlab.kordlib.common.entity.Snowflake
import com.guigou.botisgud.extensions.logger
import com.guigou.botisgud.models.User
import org.litote.kmongo.coroutine.CoroutineCollection
import org.litote.kmongo.coroutine.CoroutineDatabase
import org.litote.kmongo.eq

class UserServiceImpl(private val database: CoroutineDatabase) : UserService {
    companion object {
        val logger = logger()
    }

    private val users: CoroutineCollection<User> = database.getCollection("users")

    override suspend fun add(user: User) {
        users.insertOne(user)
    }

    override suspend fun get(userId: Snowflake): User? {
        return users.findOne(User::userId eq userId.longValue)
    }
}
