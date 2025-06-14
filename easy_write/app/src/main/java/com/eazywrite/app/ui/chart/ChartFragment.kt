package com.eazywrite.app.ui.chart

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.compose.foundation.layout.PaddingValues
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Surface
import androidx.compose.ui.platform.ComposeView
import androidx.compose.ui.unit.dp
import androidx.fragment.app.Fragment
import com.eazywrite.app.ui.theme.EazyWriteTheme

/**
 * 图表展示Fragment，使用Jetpack Compose构建UI
 * @param paddingValues 内边距设置，默认为0.dp
 */
class ChartFragment @JvmOverloads constructor(
    // 布局内边距参数，允许从外部传入自定义值
    private val paddingValues: PaddingValues = PaddingValues(0.dp)
) : Fragment() {  // 继承自Android的Fragment基类

    /**
     * 创建Fragment的视图
     * @param inflater 布局加载器
     * @param container 父容器
     * @param savedInstanceState 保存的状态数据
     * @return 创建的视图
     */
    override fun onCreateView(
        inflater: LayoutInflater,
        container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        // 创建ComposeView作为根视图
        return ComposeView(requireContext()).apply {
            // 设置Compose内容
            setContent {
                // 应用自定义主题
                EazyWriteTheme {
                    // 创建表面容器，设置背景色为主题背景色
                    Surface(color = MaterialTheme.colorScheme.background) {
                        // 显示图表页面组件，并传入paddingValues参数
                        ChartPage()
                    }
                }
            }
        }
    }
}