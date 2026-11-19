import{Ai as e,Hr as t,Jn as n,Ki as r,Mn as i,Ri as a,Yn as o,_n as s,ai as c,fi as l,ga as u,gn as d,la as f,li as p,mi as m,qi as h,ui as g,vn as _}from"./router-DhIWMPTv.js";import{n as v}from"./index-5Dh4MEKh.js";var y=p([p(`@keyframes spin-rotate`,`
 from {
 transform: rotate(0);
 }
 to {
 transform: rotate(360deg);
 }
 `),g(`spin-container`,`
 position: relative;
 `,[g(`spin-body`,`
 position: absolute;
 top: 50%;
 left: 50%;
 transform: translateX(-50%) translateY(-50%);
 `,[d()])]),g(`spin-body`,`
 display: inline-flex;
 align-items: center;
 justify-content: center;
 flex-direction: column;
 `),g(`spin`,`
 display: inline-flex;
 height: var(--n-size);
 width: var(--n-size);
 font-size: var(--n-size);
 color: var(--n-color);
 `,[l(`rotate`,`
 animation: spin-rotate 2s linear infinite;
 `)]),g(`spin-description`,`
 display: inline-block;
 font-size: var(--n-font-size);
 color: var(--n-text-color);
 transition: color .3s var(--n-bezier);
 margin-top: 8px;
 `),g(`spin-content`,`
 opacity: 1;
 transition: opacity .3s var(--n-bezier);
 pointer-events: all;
 `,[l(`spinning`,`
 user-select: none;
 -webkit-user-select: none;
 pointer-events: none;
 opacity: var(--n-opacity-spinning);
 `)])]),b={small:20,medium:18,large:16},x=r({name:`Spin`,props:Object.assign(Object.assign(Object.assign({},i.props),{contentClass:String,contentStyle:[Object,String],description:String,size:{type:[String,Number],default:`medium`},show:{type:Boolean,default:!0},rotate:{type:Boolean,default:!0},spinning:{type:Boolean,validator:()=>!0,default:void 0},delay:Number}),_),slots:Object,setup(e){let{mergedClsPrefixRef:r,inlineThemeDisabled:s}=o(e),l=i(`Spin`,`-spin`,y,v,e,r),d=a(()=>{let{size:t}=e,{common:{cubicBezierEaseInOut:n},self:r}=l.value,{opacitySpinning:i,color:a,textColor:o}=r;return{"--n-bezier":n,"--n-opacity-spinning":i,"--n-size":typeof t==`number`?c(t):r[m(`size`,t)],"--n-color":a,"--n-text-color":o}}),p=s?n(`spin`,a(()=>{let{size:t}=e;return typeof t==`number`?String(t):t[0]}),d,e):void 0,h=t(e,[`spinning`,`show`]),g=u(!1);return f(t=>{let n;if(h.value){let{delay:r}=e;if(r){n=window.setTimeout(()=>{g.value=!0},r),t(()=>{clearTimeout(n)});return}}g.value=h.value}),{mergedClsPrefix:r,active:g,mergedStrokeWidth:a(()=>{let{strokeWidth:t}=e;if(t!==void 0)return t;let{size:n}=e;return b[typeof n==`number`?`medium`:n]}),cssVars:s?void 0:d,themeClass:p?.themeClass,onRender:p?.onRender}},render(){var t;let{$slots:n,mergedClsPrefix:r,description:i}=this,a=n.icon&&this.rotate,o=(i||n.description)&&h(`div`,{class:`${r}-spin-description`},i||n.description?.call(n)),c=n.icon?h(`div`,{class:[`${r}-spin-body`,this.themeClass]},h(`div`,{class:[`${r}-spin`,a&&`${r}-spin--rotate`],style:n.default?``:this.cssVars},n.icon()),o):h(`div`,{class:[`${r}-spin-body`,this.themeClass]},h(s,{clsPrefix:r,style:n.default?``:this.cssVars,stroke:this.stroke,"stroke-width":this.mergedStrokeWidth,radius:this.radius,scale:this.scale,class:`${r}-spin`}),o);return(t=this.onRender)==null||t.call(this),n.default?h(`div`,{class:[`${r}-spin-container`,this.themeClass],style:this.cssVars},h(`div`,{class:[`${r}-spin-content`,this.active&&`${r}-spin-content--spinning`,this.contentClass],style:this.contentStyle},n),h(e,{name:`fade-in-transition`},{default:()=>this.active?c:null})):c}});export{x as t};