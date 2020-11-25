package com.guigou.botisgud

import com.gitlab.kordlib.common.entity.Snowflake
import com.gitlab.kordlib.core.Kord
import com.guigou.botisgud.commands.*
import kotlinx.coroutines.runBlocking

fun main() = runBlocking {
    val token = System.getenv("DISCORD_TOKEN")
    val client = Kord(token)

    // https://elizarov.medium.com/coroutine-context-and-scope-c8b255d59055
    client.register(Typing(), this) // TODO: determine how to initialize class within extension function
    client.register(Reminder(), this)

    val rawUsers = System.getenv("NICKNAME_USERS")
    val users = rawUsers
        .split('\n')
        .map { user ->
            user.split(':').let {
                Snowflake(it.first()) to Snowflake(it.last())
            }
        }
        .groupBy({ user -> user.first }, { user -> user.second })

    val options = System.getenv("NICKNAME_COMMAND_TRIGGER_EXPRESSION").let {
        if (it.isNullOrEmpty()) NicknameOptions() else NicknameOptions(it)
    }
    client.register(Nickname(users, options), this)

    client.login()
}
