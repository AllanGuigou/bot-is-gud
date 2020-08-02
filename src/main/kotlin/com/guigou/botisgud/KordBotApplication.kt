package com.guigou.botisgud

import com.gitlab.kordlib.core.Kord
import com.guigou.botisgud.commands.Reminder
import com.guigou.botisgud.commands.Typing
import com.guigou.botisgud.commands.register

suspend fun main() {
    val client = Kord(System.getenv("DISCORD_TOKEN"))

    client.register(Typing()) // TODO: determine how to initialize class within extension function
    client.register(Reminder())

    client.login()
}
