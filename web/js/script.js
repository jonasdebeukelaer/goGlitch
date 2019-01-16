Vue.component('effect-card', {
    props: ['name', 'id'],
    template: `
    <el-card class="effect-card">
        <div slot="header" class="clearfix header">
            <span>{{ name }}</span>
            <el-button class="del-effect-btn">
                <i class="el-icon-close"></i>
            </el-button>
        </div>
        <div class="item effect-control">
            adjustments will go here
        </div>
    </el-card>
    `
});

var app = new Vue({
    el: '#app',
    data: {
        uploadingImage: false,
        processingImage: false,
        introTextClass: '',
        mainImageClass: 'hidden',
        processBtnText: 'process image',
        sidebarVisibleClass: 'hidden',
        spinnerVisibleClass: 'hidden',
        filename: '',
        sidebarImgUrl: 'source_image/placeholder.png',
        mainImgUrl: 'source_image/placeholder.png',
        effectOptions: [],
        effectLayers: []
    },
    methods: {
        showInfo() {
            this.$alert('welcome to goGlitch, where you can use some tools I\'ve created for glitching images', 
            'Welcome to goglitch', {
                confirmButtonText: 'OK'
              });
        },
        selectFile() {
            this.$refs.image_file.click()
        },
        submitImage(){
            var filelist = this.$refs.image_file.files
            if (filelist.length !== 0) {
                this.image_file = filelist[0];
                uploadImage(this)
            }
        },
        processImage() {
            processAndLoadImage(this)
        },
        newEffectLayer(selectedEffect) {
            addNewEffectLayer(this, selectedEffect)
        }
    }
})

var addNewEffectLayer = function(that, selectedEffect) {
    that.effectLayers.push({
        name: selectedEffect,
        key: selectedEffect,
        id: selectedEffect + that.effectLayers.length
    })
}

var uploadImage = function(that) {
    that.uploadingImage = true;
    let formData = new FormData();
    formData.append("image", that.image_file);
    axios.post( '/upload_image',
        formData,
        {
            headers: {'Content-Type': 'multipart/form-data'}
        }
    ).then(function(res){
        that.uploadingImage = false;
        that.filename = res.data.filename;
        that.sidebarImgUrl = "source_image/" + res.data.filename;
        that.mainImgUrl = "source_image/placeholder.png" ;
        that.sidebarVisibleClass = '';
        that.introTextClass = 'hidden';
        
        getEffectOptions(that, function(effects) {
            console.log(effects)
            if (effects !== undefined && effects.length !== 0) {
                effects.forEach(e => {
                    that.effectOptions.push(e);
                }); 
            }
        })
        
    })
    .catch(function(res){
        reportError(that, 'failed to upload image!');
        console.log(res);
    }).then(function() {
        that.uploadingImage = false;
    })
}

var getEffectOptions = function(that, callback) {
    axios.get( '/effect_options',
    ).then(function(res){
        callback(res.data);
    })
    .catch(function(res){
        reportError(that, 'failed to get effect options list!');
        console.log(res);
    });
}

var processAndLoadImage = function(that) {
    that.processBtnText = '';
    that.processingImage = true;
    that.spinnerVisibleClass = '';

    var data = [];
    that.effectLayers.forEach(e => {
        return data.push({"key":e.key});
    })

    dataStr = JSON.stringify(data);
    console.log(dataStr)
    var url = '/process_image?image=' + that.filename;
    axios.post( url,
        dataStr,
        {
            headers: {'Content-Type': 'application/json'}
        }
        ).then(function(res){
        that.mainImgUrl = "processed_image/" + res.data.filename;
        that.spinnerVisibleClass = 'hidden';
        that.mainImageClass = '';
        console.log('image processed');
    })
    .catch(function(){
        reportError(that, 'failed to process image!');
    })
    .then(function() {
        that.processingImage = false;
        that.processBtnText = 'process image';
    });
}

var reportError = function(that, errMsg) {
    const h = that.$createElement;

    that.$notify({
      title: 'Error',
      message: h('i', { style: 'color: #202020' }, errMsg)
    });
}