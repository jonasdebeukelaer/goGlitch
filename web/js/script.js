var app = new Vue({
    el: '#app',
    data: {
        message: 'Hello Vue!',
        workspaceVisibleClass: 'hidden',
        spinnerVisibleClass: 'hidden',
        filename: '',
        sourceImgUrl: 'source_image/placeholder.png',
        processedImgUrl: 'source_image/placeholder.png'
    },
    methods: {
        handleImageUpload() {
            this.image_file = this.$refs.image_file.files[0];
        },
        submitImage(){
            //showSpinner()
            uploadImage(this, function() {
            })
            //hideSpinner()
        },
        processImage() {
            processAndLoadImage(this)
        }
    }
})

var uploadImage = function(that, callback) {
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
            that.filename = res.data.filename
            that.sourceImgUrl = "source_image/" + res.data.filename
            that.processedImgUrl = "source_image/placeholder.png" 
            that.workspaceVisibleClass = ''
            console.log('image uploaded');
        })
        .catch(function(){
            console.log('failed to upload image');
    });
}

var processAndLoadImage = function(that) {
    that.spinnerVisibleClass = ''
    url = '/process_image?image=' + that.filename
    axios.get( url ).then(function(res){
            that.processedImgUrl = "processed_image/" + res.data.filename
            that.spinnerVisibleClass = 'hidden'
            console.log('image processed');
        })
        .catch(function(){
            console.log('failed to upload image');
    });
}
