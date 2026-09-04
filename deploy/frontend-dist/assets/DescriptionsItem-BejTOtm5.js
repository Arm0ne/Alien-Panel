import{$r as e,Cr as t,Fi as n,Jr as r,Kn as i,Li as a,Oi as o,Pr as s,Qr as c,Xr as l,Yr as u,Zr as d,ei as f,ir as p,kn as m,qn as h,qr as g}from"./router-Ca5XLGLk.js";import{n as _}from"./Space-BRTVS2IW.js";import{s as v}from"./index-ChrUCjfM.js";function y(e,t=`default`,n=[]){let{children:r}=e;if(typeof r==`object`&&r&&!Array.isArray(r)){let e=r[t];if(typeof e==`function`)return e()}return n}var b=g([r(`descriptions`,{fontSize:`var(--n-font-size)`},[r(`descriptions-separator`,`
 display: inline-block;
 margin: 0 8px 0 2px;
 `),r(`descriptions-table-wrapper`,[r(`descriptions-table`,[r(`descriptions-table-row`,[r(`descriptions-table-header`,{padding:`var(--n-th-padding)`}),r(`descriptions-table-content`,{padding:`var(--n-td-padding)`})])])]),d(`bordered`,[r(`descriptions-table-wrapper`,[r(`descriptions-table`,[r(`descriptions-table-row`,[g(`&:last-child`,[r(`descriptions-table-content`,{paddingBottom:0})])])])])]),l(`left-label-placement`,[r(`descriptions-table-content`,[g(`> *`,{verticalAlign:`top`})])]),l(`left-label-align`,[g(`th`,{textAlign:`left`})]),l(`center-label-align`,[g(`th`,{textAlign:`center`})]),l(`right-label-align`,[g(`th`,{textAlign:`right`})]),l(`bordered`,[r(`descriptions-table-wrapper`,`
 border-radius: var(--n-border-radius);
 overflow: hidden;
 background: var(--n-merged-td-color);
 border: 1px solid var(--n-merged-border-color);
 `,[r(`descriptions-table`,[r(`descriptions-table-row`,[g(`&:not(:last-child)`,[r(`descriptions-table-content`,{borderBottom:`1px solid var(--n-merged-border-color)`}),r(`descriptions-table-header`,{borderBottom:`1px solid var(--n-merged-border-color)`})]),r(`descriptions-table-header`,`
 font-weight: 400;
 background-clip: padding-box;
 background-color: var(--n-merged-th-color);
 `,[g(`&:not(:last-child)`,{borderRight:`1px solid var(--n-merged-border-color)`})]),r(`descriptions-table-content`,[g(`&:not(:last-child)`,{borderRight:`1px solid var(--n-merged-border-color)`})])])])])]),r(`descriptions-header`,`
 font-weight: var(--n-th-font-weight);
 font-size: 18px;
 transition: color .3s var(--n-bezier);
 line-height: var(--n-line-height);
 margin-bottom: 16px;
 color: var(--n-title-text-color);
 `),r(`descriptions-table-wrapper`,`
 transition:
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `,[r(`descriptions-table`,`
 width: 100%;
 border-collapse: separate;
 border-spacing: 0;
 box-sizing: border-box;
 `,[r(`descriptions-table-row`,`
 box-sizing: border-box;
 transition: border-color .3s var(--n-bezier);
 `,[r(`descriptions-table-header`,`
 font-weight: var(--n-th-font-weight);
 line-height: var(--n-line-height);
 display: table-cell;
 box-sizing: border-box;
 color: var(--n-th-text-color);
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `),r(`descriptions-table-content`,`
 vertical-align: top;
 line-height: var(--n-line-height);
 display: table-cell;
 box-sizing: border-box;
 color: var(--n-td-text-color);
 transition:
 color .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 border-color .3s var(--n-bezier);
 `,[u(`content`,`
 transition: color .3s var(--n-bezier);
 display: inline-block;
 color: var(--n-td-text-color);
 `)]),u(`label`,`
 font-weight: var(--n-th-font-weight);
 transition: color .3s var(--n-bezier);
 display: inline-block;
 margin-right: 14px;
 color: var(--n-th-text-color);
 `)])])])]),r(`descriptions-table-wrapper`,`
 --n-merged-th-color: var(--n-th-color);
 --n-merged-td-color: var(--n-td-color);
 --n-merged-border-color: var(--n-border-color);
 `),e(r(`descriptions-table-wrapper`,`
 --n-merged-th-color: var(--n-th-color-modal);
 --n-merged-td-color: var(--n-td-color-modal);
 --n-merged-border-color: var(--n-border-color-modal);
 `)),f(r(`descriptions-table-wrapper`,`
 --n-merged-th-color: var(--n-th-color-popover);
 --n-merged-td-color: var(--n-td-color-popover);
 --n-merged-border-color: var(--n-border-color-popover);
 `))]),x=`DESCRIPTION_ITEM_FLAG`;function S(e){return typeof e==`object`&&e&&!Array.isArray(e)?e.type&&e.type.DESCRIPTION_ITEM_FLAG:!1}var C=n({name:`Descriptions`,props:Object.assign(Object.assign({},m.props),{title:String,column:{type:Number,default:3},columns:Number,labelPlacement:{type:String,default:`top`},labelAlign:{type:String,default:`left`},separator:{type:String,default:`:`},size:String,bordered:Boolean,labelClass:String,labelStyle:[Object,String],contentClass:String,contentStyle:[Object,String]}),slots:Object,setup(e){let{mergedClsPrefixRef:n,inlineThemeDisabled:r,mergedComponentPropsRef:a}=h(e),s=o(()=>e.size||a?.value?.Descriptions?.size||`medium`),l=m(`Descriptions`,`-descriptions`,b,v,e,n),u=o(()=>{let{bordered:t}=e,n=s.value,{common:{cubicBezierEaseInOut:r},self:{titleTextColor:i,thColor:a,thColorModal:o,thColorPopover:u,thTextColor:d,thFontWeight:f,tdTextColor:p,tdColor:m,tdColorModal:h,tdColorPopover:g,borderColor:_,borderColorModal:v,borderColorPopover:y,borderRadius:b,lineHeight:x,[c(`fontSize`,n)]:S,[c(t?`thPaddingBordered`:`thPadding`,n)]:C,[c(t?`tdPaddingBordered`:`tdPadding`,n)]:w}}=l.value;return{"--n-title-text-color":i,"--n-th-padding":C,"--n-td-padding":w,"--n-font-size":S,"--n-bezier":r,"--n-th-font-weight":f,"--n-line-height":x,"--n-th-text-color":d,"--n-td-text-color":p,"--n-th-color":a,"--n-th-color-modal":o,"--n-th-color-popover":u,"--n-td-color":m,"--n-td-color-modal":h,"--n-td-color-popover":g,"--n-border-radius":b,"--n-border-color":_,"--n-border-color-modal":v,"--n-border-color-popover":y}}),d=r?i(`descriptions`,o(()=>{let t=``,{bordered:n}=e;return n&&(t+=`a`),t+=s.value[0],t}),u,e):void 0;return{mergedClsPrefix:n,cssVars:r?void 0:u,themeClass:d?.themeClass,onRender:d?.onRender,compitableColumn:t(e,[`columns`,`column`]),inlineThemeDisabled:r,mergedSize:s}},render(){let e=this.$slots.default,t=e?p(e()):[];t.length;let{contentClass:n,labelClass:r,compitableColumn:i,labelPlacement:o,labelAlign:c,mergedSize:l,bordered:u,title:d,cssVars:f,mergedClsPrefix:m,separator:h,onRender:g}=this;g?.();let v=t.filter(e=>S(e)),b=v.reduce((e,t,s)=>{let c=t.props||{},l=v.length-1===s,d=[`label`in c?c.label:y(t,`label`)],f=[y(t)],p=c.span||1,g=e.span;e.span+=p;let _=c.labelStyle||c[`label-style`]||this.labelStyle,b=c.contentStyle||c[`content-style`]||this.contentStyle;if(o===`left`)u?e.row.push(a(`th`,{class:[`${m}-descriptions-table-header`,r],colspan:1,style:_},d),a(`td`,{class:[`${m}-descriptions-table-content`,n],colspan:l?(i-g)*2+1:p*2-1,style:b},f)):e.row.push(a(`td`,{class:`${m}-descriptions-table-content`,colspan:l?(i-g)*2:p*2},a(`span`,{class:[`${m}-descriptions-table-content__label`,r],style:_},[...d,h&&a(`span`,{class:`${m}-descriptions-separator`},h)]),a(`span`,{class:[`${m}-descriptions-table-content__content`,n],style:b},f)));else{let t=l?(i-g)*2:p*2;e.row.push(a(`th`,{class:[`${m}-descriptions-table-header`,r],colspan:t,style:_},d)),e.secondRow.push(a(`td`,{class:[`${m}-descriptions-table-content`,n],colspan:t,style:b},f))}return(e.span>=i||l)&&(e.span=0,e.row.length&&(e.rows.push(e.row),e.row=[]),o!==`left`&&e.secondRow.length&&(e.rows.push(e.secondRow),e.secondRow=[])),e},{span:0,row:[],secondRow:[],rows:[]}).rows.map(e=>a(`tr`,{class:`${m}-descriptions-table-row`},e));return a(`div`,{style:f,class:[`${m}-descriptions`,this.themeClass,`${m}-descriptions--${o}-label-placement`,`${m}-descriptions--${c}-label-align`,`${m}-descriptions--${l}-size`,u&&`${m}-descriptions--bordered`]},d||this.$slots.header?a(`div`,{class:`${m}-descriptions-header`},d||_(this,`header`)):null,a(`div`,{class:`${m}-descriptions-table-wrapper`},a(`table`,{class:`${m}-descriptions-table`},a(`tbody`,null,o===`top`&&a(`tr`,{class:`${m}-descriptions-table-row`,style:{visibility:`collapse`}},s(i*2,a(`td`,null))),b))))}}),w={label:String,span:{type:Number,default:1},labelClass:String,labelStyle:[Object,String],contentClass:String,contentStyle:[Object,String]},T=n({name:`DescriptionsItem`,[x]:!0,props:w,slots:Object,render(){return null}});export{C as n,T as t};