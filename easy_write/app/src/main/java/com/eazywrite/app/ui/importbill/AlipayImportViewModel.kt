package com.eazywrite.app.ui.importbill

import android.app.Application
import com.eazywrite.app.data.model.ImportBillData
import com.github.doyaaaaaken.kotlincsv.dsl.context.ExcessFieldsRowBehaviour
import com.github.doyaaaaaken.kotlincsv.dsl.context.InsufficientFieldsRowBehaviour
import com.github.doyaaaaaken.kotlincsv.dsl.csvReader
import java.io.File
import java.time.LocalDateTime
import java.time.format.DateTimeFormatter

/**
 * 支付宝账单导入ViewModel
 * 负责解析支付宝导出的CSV账单文件并将其转换为应用内部数据模型
 */
class AlipayImportViewModel(application: Application) : ImportViewModel(application) {

    /**
     * 实现抽象导入方法
     * @param csvFile 支付宝导出的CSV文件
     * @return 解析后的账单数据列表
     */
    override suspend fun importImpl(csvFile: File) = importAlipay(csvFile)

    /**
     * 解析支付宝CSV账单文件的核心方法
     * @param csvFile 支付宝导出的CSV文件
     * @return 转换后的账单数据列表
     */
    private suspend fun importAlipay(csvFile: File): List<ImportBillData> {
        // 配置CSV阅读器
        val bills = csvReader {
            // 处理字段不足的行：填充空字符串
            insufficientFieldsRowBehaviour = InsufficientFieldsRowBehaviour.EMPTY_STRING
            // 处理多余字段的行：自动截断
            excessFieldsRowBehaviour = ExcessFieldsRowBehaviour.TRIM
            // 支付宝账单使用GBK编码
            charset = "GBK"
        }.openAsync(csvFile) {
            val bills = mutableListOf<ImportBillData>()
            var startIndex = -1  // 记录实际数据开始的行索引

            // 从第13行开始读取（跳过支付宝的文件头信息）
            readAllAsSequence(13).forEachIndexed { index, row0 ->
                // 清理每行的数据：去除首尾空格
                val row = row0.map { it.trim() }

                // 检测数据起始标记："电子客户回单"
                if (row[0].contains("电子客户回单")) {
                    startIndex = index + 2  // 数据在标记行后2行开始
                }

                // 处理有效数据行
                if (startIndex != -1 && index >= startIndex) {
                    // 创建日期时间解析器（支付宝格式：yyyy-MM-dd HH:mm:ss）
                    val dateTimeFormatter = DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss")

                    // 解析交易时间
                    val datetime = LocalDateTime.from(dateTimeFormatter.parse(row[0]))

                    // 转换交易类型：支付宝的"收入/支出"转成内部标识
                    val type = when (row[5]) {
                        "收入" -> "in"
                        "支出" -> "out"
                        else -> "other"  // 其他类型（如退款）暂时跳过
                    }

                    // 跳过非收支类型的交易
                    if (type == "other") return@forEachIndexed

                    // 解析交易金额（第6列）
                    val amount = row[6].toBigDecimal()

                    // 处理备注信息（第11列）
                    val comment = if (row[11] == "/") "" else row[11]

                    // 处理交易方名称（第4列）
                    val name = if (row[4] == "/") "未命名" else row[4]

                    // 构建账单数据对象
                    bills.add(
                        ImportBillData(
                            datetime = datetime,         // 交易时间
                            category = row[1],            // 交易分类
                            transactionPartner = row[2],  // 交易对方
                            partnerAccount = row[3],      // 对方账号
                            name = name,                  // 商品名称
                            type = type,                  // 交易类型（in/out）
                            amount = amount,              // 交易金额
                            paymentMethod = row[7],       // 收/付款方式
                            status = row[8],              // 交易状态
                            transactionNo = "alipay-" + row[9],  // 交易流水号（添加支付宝前缀）
                            shopNo = row[10],             // 商户订单号
                            comment = comment             // 交易备注
                        )
                    )
                }
            }
            return@openAsync bills.toList()  // 返回最终解析结果
        }
        return bills
    }
}