package com.eazywrite.app.common

import androidx.fragment.app.FragmentManager
import com.eazywrite.app.R
import com.google.android.material.datepicker.MaterialDatePicker
import com.google.android.material.timepicker.MaterialTimePicker
import com.google.android.material.timepicker.TimeFormat

/**
 * 弹出一个日期选择器供用户选择日期。
 *
 * @param manager FragmentManager对象，用于显示日期选择器。
 * @param currentTimestamp 当前选择的日期时间戳，用于设置日期选择器的初始日期。
 * @param onDateSelected 当用户选择日期后调用的回调函数。
 */
fun pickDate(manager: FragmentManager, currentTimestamp: Long, onDateSelected: (Long) -> Unit) {

    // 创建日期选择器的构建器
    val datePicker =
        MaterialDatePicker.Builder.datePicker()
            .setTitleText("请选择日期") // 设置日期选择器的标题
            .setTheme(R.style.ThemeOverlay_App_DatePicker) // 设置日期选择器的主题样式
            .setSelection(currentTimestamp) // 设置日期选择器的初始选择为 currentTimestamp
            .build() // 构建日期选择器
            .apply {
                addOnPositiveButtonClickListener { timestamp ->
                    onDateSelected(timestamp) // 用户选择日期后调用 onDateSelected 回调
                }
            }
    datePicker.show(manager, "datePicker") // 在FragmentManager中显示日期选择器
}

/**
 * 弹出一个时间选择器供用户选择时间。
 *
 * @param manager FragmentManager对象，用于显示时间选择器。
 * @param hour 默认小时，用于设置时间选择器的初始小时。
 * @param minute 默认分钟，用于设置时间选择器的初始分钟。
 * @param onDateSelected 当用户选择时间后调用的回调函数。
 */
fun pickTime(
    manager: FragmentManager,
    hour: Int = 12, // 默认小时设置为12
    minute: Int = 0, // 默认分钟设置为0
    onDateSelected: (hour: Int, minute: Int) -> Unit // 用户选择时间后调用的回调函数
) {
    // 创建时间选择器的构建器
    val picker =
        MaterialTimePicker.Builder()
            .setTimeFormat(TimeFormat.CLOCK_24H) // 设置时间格式为24小时制
            .setHour(hour) // 设置时间选择器的初始小时
            .setMinute(minute) // 设置时间选择器的初始分钟
            .setTitleText("请选择时间") // 设置时间选择器的标题
            .build() // 构建时间选择器

    // 为时间选择器的确认按钮设置点击监听器
    picker.addOnPositiveButtonClickListener {
        onDateSelected(picker.hour, picker.minute) // 用户选择时间后调用 onDateSelected 回调
    }

    // 在FragmentManager中显示时间选择器
    picker.show(manager, "timePicker")
}