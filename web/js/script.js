window.Event = new Vue();

Vue.component('effect-card', {
    data: {
        params: {}
    },
    props: ['name', 'effect_key', 'id'],
    methods: {
        removeEffect(){
            Event.$emit('removeEffect', this.id)
       }
    },
    template: `
    <el-card class="effect-card">
        <div slot="header" class="clearfix header">
            <span>{{ name }}</span>
            <el-button circle :id="'effect'+id" class="del-effect-btn" @click="removeEffect()">
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
        imageUploaded: false,
        imageProcessed: false,
        processBtnText: 'process image',
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
            this.$refs.image_file.click();
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
        },
        removeEffect(effectId) {this.effectLayers.pop(effectId);}
    },
    created(){
        Event.$on('removeEffect', this.removeEffect)
  }
})

var addNewEffectLayer = function(that, selectedEffect) {
    that.effectLayers.push({
        name: selectedEffect.name,
        effect_key: selectedEffect.effect_key,
        id: that.effectLayers.length,
        params: selectedEffect.params
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
        that.imageUploaded = true;
        that.introTextClass = 'hidden';
        
        while (that.effectOptions.length > 0) {
            that.effectOptions.pop()
        }
        getEffectOptions(that, function(effects) {
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
        console.log(res.data)
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

    var data = [];
    console.log(that.effectLayers)
    that.effectLayers.forEach(e => {
        return data.push({
            "key":e.effect_key,
            "params": {
                "offsetPercent": "0.5"
            }
        });
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
        that.mainImageClass = '';
        console.log('image processed');
    })
    .catch(function(){
        reportError(that, 'failed to process image!');
    })
    .then(function() {
        that.processingImage = false;
        that.imageProcessed = true;
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