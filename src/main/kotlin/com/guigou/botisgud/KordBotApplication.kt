package com.guigou.botisgud

import com.gitlab.kordlib.core.Kord
import com.guigou.botisgud.commands.register
import com.guigou.botisgud.services.nickname.Nicknames
import com.guigou.botisgud.services.reminder.Reminders
import com.guigou.botisgud.services.user.Users
import kotlinx.coroutines.runBlocking
import org.jetbrains.exposed.sql.Database
import org.jetbrains.exposed.sql.SchemaUtils
import org.jetbrains.exposed.sql.transactions.transaction

fun main() = runBlocking {
    val token = System.getenv("DISCORD_TOKEN")
    val client = Kord(token)

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

    val component = ApplicationComponent::class.create()
    client.register(component.nicknameCommand, this)
    client.register(component.reminderCommand, this)
    client.register(component.typingCommand, this)
    client.login()
}
