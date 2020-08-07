package com.guigou.botisgud.extensions.kord

import com.gitlab.kordlib.core.event.message.ReactionAddEvent

val ReactionAddEvent.link: String
    get() = "https://discordapp.com/channels/${guildId?.value ?: "@me"}/${channelId.value}/${messageId.value}"
