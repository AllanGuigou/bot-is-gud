package com.guigou.botisgud.commands

import com.gitlab.kordlib.core.Kord

interface Command {
    fun register(client: Kord)
}

fun <T : Command> Kord.register(command: T) {
    command.register(this)
}
