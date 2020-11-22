package com.guigou.botisgud

import com.gitlab.kordlib.core.Kord
import com.guigou.botisgud.commands.Nicknamer
import com.guigou.botisgud.commands.Reminder
import com.guigou.botisgud.commands.Typing
import com.guigou.botisgud.commands.register

suspend fun main() {
    val token = System.getenv("DISCORD_TOKEN")
    val client = Kord(token)

    client.register(Typing()) // TODO: determine how to initialize class within extension function
    client.register(Reminder())
    client.register(Nicknamer())

    client.login()
}
