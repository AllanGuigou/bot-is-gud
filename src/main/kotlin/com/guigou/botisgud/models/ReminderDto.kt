package com.guigou.botisgud.models

import com.gitlab.kordlib.common.entity.Snowflake

// TODO: make link a uri
data class ReminderDto(val userId: Snowflake, val message: String, val link: String)
