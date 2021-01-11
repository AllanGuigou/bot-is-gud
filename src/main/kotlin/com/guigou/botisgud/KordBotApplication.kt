package com.guigou.botisgud

import com.gitlab.kordlib.common.entity.Snowflake
import com.gitlab.kordlib.core.Kord
import com.guigou.botisgud.commands.*
import com.guigou.botisgud.services.nickname.NicknameEntity
import com.guigou.botisgud.services.nickname.Nicknames
import com.guigou.botisgud.services.reminder.ReminderServiceImpl
import com.guigou.botisgud.services.reminder.Reminders
import com.guigou.botisgud.services.user.Users
import kotlinx.coroutines.runBlocking
import org.jetbrains.exposed.sql.*
import org.jetbrains.exposed.sql.transactions.transaction

fun main() = runBlocking {
    val token = System.getenv("DISCORD_TOKEN")
    val client = Kord(token)

    // https://elizarov.medium.com/coroutine-context-and-scope-c8b255d59055
    client.register(Typing(), this) // TODO: determine how to initialize class within extension function

    Database.connect(
        System.getenv("DATABASE_URL"),
        System.getenv("DATABASE_DRIVER"),
        System.getenv("DATABASE_USER"),
        System.getenv("DATABASE_PASSWORD")
    )

    // TODO: determine how a migration tool to use
    transaction {
        SchemaUtils.create(Users)
        SchemaUtils.create(Reminders)
        SchemaUtils.create(Nicknames)
    }

    client.register(Reminder(ReminderServiceImpl()), this)

    val options = System.getenv("NICKNAME_COMMAND_TRIGGER_EXPRESSION").let {
        if (it.isNullOrEmpty()) NicknameOptions() else NicknameOptions(it)
    }

    client.register(Nickname(options), this)

    client.login()
}
