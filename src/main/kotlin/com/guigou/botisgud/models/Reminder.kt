package com.guigou.botisgud.models

import com.gitlab.kordlib.common.entity.Snowflake
import java.util.*

data class Reminder(val userId: Snowflake, val message: String, val link: String, val timestamp: Date)
