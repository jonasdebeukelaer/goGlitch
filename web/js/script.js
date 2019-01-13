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
        fileList: []
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
        reportError(errMsg) {
            const h = this.$createElement;
    
            this.$notify({
              title: 'Error',
              message: h('i', { style: 'color: #202020' }, errMsg)
            });
          },
    }
})

var uploadImage = function(that) {
    that.uploadingImage = true
    let formData = new FormData();
    formData.append("image", that.image_file)
    axios.post( '/upload_image',
        formData,
        {
            headers: {
                'Content-Type': 'multipart/form-data'
            }
        }
        ).then(function(res){
            that.uploadingImage = false
            that.filename = res.data.filename
            that.sidebarImgUrl = "source_image/" + res.data.filename
            that.mainImgUrl = "source_image/placeholder.png" 
            that.sidebarVisibleClass = ''
            that.introTextClass = 'hidden'
        })
        .catch(function(res){
            reportError('failed to upload image!');
            console.log(res)
        .then(function() {
            that.uploadingImage = false
        })
    });
}

var processAndLoadImage = function(that) {
    that.processBtnText = ''
    that.processingImage = true
    that.spinnerVisibleClass = ''
    url = '/process_image?image=' + that.filename
    axios.get( url ).then(function(res){
        that.mainImgUrl = "processed_image/" + res.data.filename
        that.spinnerVisibleClass = 'hidden'
        that.mainImageClass = ''
        console.log('image processed');
    })
    .catch(function(){
        reportError('failed to process image!');
    })
    .then(function() {
        that.processingImage = false
        that.processBtnText = 'process image'
    });
}
