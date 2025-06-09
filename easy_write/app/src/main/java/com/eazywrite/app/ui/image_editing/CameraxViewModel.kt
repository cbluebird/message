package com.eazywrite.app.ui.image_editing

import android.app.Application
import androidx.compose.runtime.mutableStateListOf
import androidx.lifecycle.AndroidViewModel
import java.io.File


class CameraxViewModel(application: Application) : AndroidViewModel(application) {

    val images = mutableStateListOf<File>()

    override fun onCleared() {
        super.onCleared()
        images.forEach {
            it.delete()
        }
    }
}